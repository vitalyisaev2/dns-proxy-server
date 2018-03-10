package docker

import (
	"fmt"
	"github.com/docker/engine-api/types/filters"
	"encoding/json"
	"github.com/docker/engine-api/types/events"
	"io"
	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
	log "github.com/mageddo/go-logging"
	"github.com/docker/engine-api/types"
	"github.com/mageddo/dns-proxy-server/cache/lru"
	"github.com/mageddo/dns-proxy-server/cache"
	"strings"
	"errors"
	"github.com/mageddo/dns-proxy-server/conf"
)

var c = lru.New(43690)

func HandleDockerEvents(){
	defaultLogger := log.NewContext()
	logger := log.NewLog(defaultLogger)

	// connecting to docker api
	cli, err := client.NewClient(conf.DockerHost(), "v1.21", nil, nil)
	if err != nil {
		logger.Errorf("status=error-to-connect-at-host, solver=docker, err=%v", err)
		return
	}

	// more about list containers https://docs.docker.com/engine/reference/commandline/ps/
	options := types.ContainerListOptions{}
	ctx := context.Background()

	serverVersion, err := cli.ServerVersion(ctx)
	logger.Infof("serverVersion=%+v, err=%v", serverVersion, err)

	containers, err := cli.ContainerList(ctx, options)
	if err != nil {
		logger.Errorf("status=error-to-list-container, solver=docker, err=%v", err)
		return
	}

	// more about events here: http://docs-stage.docker.com/v1.10/engine/reference/commandline/events/
	var eventFilter = filters.NewArgs()
	eventFilter.Add("event", "start")
	eventFilter.Add("event", "die")
	eventFilter.Add("event", "stop")

	// registering at events before get the list of actual containers, this way no one container will be missed #55
	body, err := cli.Events(ctx, types.EventsOptions{ Filters: eventFilter })
	if err != nil {
		logger.Errorf("status=error-to-attach-at-events-handler, solver=docker, err=%v", err)
		return
	}

	for _, c := range containers {

		cInspection, err := cli.ContainerInspect(ctx, c.ID)
		logger.Infof("status=container-from-list-begin, container=%s", cInspection.Name)
		if err != nil {
			logger.Errorf("status=inspect-error-at-list, container=%s, err=%v", c.Names, err)
		}
		hostnames := getHostnames(cInspection)
		putHostnames(defaultLogger, hostnames, cInspection)
		logger.Infof("status=container-from-list-success, container=%s, hostnames=%s", cInspection.Name, hostnames)

	}
	
	dec := json.NewDecoder(body)
	for {

		logCtx := log.NewContext()
		logger = log.NewLog(logCtx)

		var event events.Message
		err := dec.Decode(&event)
		if err != nil && err == io.EOF {
			break
		}

		cInspection, err := cli.ContainerInspect(ctx, event.ID)
		if err != nil {
			logger.Warningf("status=inspect-error, id=%s, err=%v", event.ID, err)
			continue
		}
		hostnames := getHostnames(cInspection)
		action := event.Action
		if len(action) == 0 {
			action = event.Status
		}
		logger.Infof("status=resolved-hosts, action=%s, hostnames=%s", action, hostnames)

		switch action {
		case "start":
			putHostnames(logCtx, hostnames, cInspection)
			break

		case "stop", "die":
			for _, host := range hostnames {
				c.Remove(host)
			}
			break

		}
	}

}

func GetCache() cache.Cache {
	return c
}

func getHostnames(inspect types.ContainerJSON) []string {

	const hostnameEnv = "HOSTNAMES="
	hostnames := *new([]string)

	if len(inspect.Config.Hostname) != 0 {
		if len(inspect.Config.Domainname) != 0 {
			hostnames = append(hostnames, fmt.Sprintf("%s.%s", inspect.Config.Hostname, inspect.Config.Domainname))
		}else {
			hostnames = append(hostnames, inspect.Config.Hostname)
		}
	}
	for _, env := range inspect.Config.Env {

		envName := strings.Index(env, hostnameEnv)
		if envName == 0 {
			envValue := env[envName + len(hostnameEnv) : ]
			hostnames = append(hostnames, strings.Split(envValue, ",")...)
			return hostnames
		}
	}
	return hostnames
}

func putHostnames(ctx context.Context, hostnames []string, inspect types.ContainerJSON) error {
	logger := log.NewLog(ctx)
	for _, host := range hostnames {

		var ip string = ""
		for _, network := range inspect.NetworkSettings.Networks {
			ip = network.IPAddress
		}
		if len(ip) == 0 {
			ip = inspect.NetworkSettings.IPAddress;
			if len(ip) == 0 {
				err := fmt.Sprintf("no network found to %s", inspect.Name)
				logger.Error(err)
				return errors.New(err)
			}
		}
		logger.Debugf("host=%s, ip=%s", host, ip)
		c.Put(host, ip)
	}
	return nil
}
