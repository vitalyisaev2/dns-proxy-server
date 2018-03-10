package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"os"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"github.com/mageddo/go-logging"
	"strings"
	. "github.com/mageddo/dns-proxy-server/log"
)

func CpuProfile() string {
	return *flags.Cpuprofile
}

func Compress() bool {
	return *flags.Compress
}

func Tsig() string {
	return *flags.Tsig
}

func WebServerPort() int {
	if conf, _ := getConf(); conf != nil && conf.WebServerPort > 0 {
		return conf.WebServerPort
	}
	return *flags.WebServerPort
}

func DnsServerPort() int {
	if conf, _ := getConf(); conf != nil && conf.DnsServerPort > 0 {
		return conf.DnsServerPort
	}
	return *flags.DnsServerPort
}

func SetupResolvConf() bool {
	if conf, _ := getConf(); conf != nil && conf.DefaultDns != nil {
		LOGGER.Debugf("status=from-conf")
		return *conf.DefaultDns
	}
	LOGGER.Debugf("status=from-flag, value=%t", *flags.SetupResolvconf)
	return *flags.SetupResolvconf
}

func GetResolvConf() string {
	return GetString(os.Getenv(env.RESOLVCONF), "/etc/resolv.conf")
}

func getConf() (*local.LocalConfiguration, error) {
	return local.LoadConfiguration(logging.NewContext())
}

func LogLevel() string {
	if lvl := os.Getenv(env.LOG_LEVEL); lvl != "" {
		return lvl
	}

	if conf, _ := getConf(); conf != nil && conf.LogLevel != "" {
		return conf.LogLevel
	}
	return flags.LogLevel()
}

func LogFile() string {
	f := os.Getenv(env.LOG_FILE)
	if conf, _ := getConf(); f == "" &&  conf != nil && conf.LogFile != "" {
		f = conf.LogFile
	}
	f = GetString(f, flags.LogToFile())

	if strings.ToLower(f) == "true" {
		return "/var/log/dns-proxy-server.log"
	}
	if strings.ToLower(f) == "false" {
		return ""
	}
	return f
}

func DockerHost() string {
	return GetString(os.Getenv(env.DOCKER_HOST), *flags.DockerHost)
}

func GetString(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
