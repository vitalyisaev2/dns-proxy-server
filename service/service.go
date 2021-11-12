package service

import (
    "fmt"
    "github.com/mageddo/dns-proxy-server/flags"
    "github.com/mageddo/dns-proxy-server/utils"
    "os"
    "errors"
    "github.com/mageddo/go-logging"
)

const (
    DNS_PROXY_SERVER_SERVICE = "dns-proxy-server"
    DNS_PROXY_SERVER_PATH = "/etc/init.d/dns-proxy-server"
)

type Service struct {}

type Script struct {
    Script string
}

func NewService() *Service {
    return &Service{}
}

func (sc *Service) Install() {

    serviceMode := *flags.SetupService
    if len(serviceMode) == 0 {
        return
    }
    logging.Infof("mode=%s, version=%s", serviceMode, flags.GetRawCurrentVersion())
    var err error
    switch serviceMode {
    case "docker":
        err = sc.SetupFor(DNS_PROXY_SERVER_PATH, DNS_PROXY_SERVER_SERVICE, NewDockerScript())
    case "normal":
        err = sc.SetupFor(DNS_PROXY_SERVER_PATH, DNS_PROXY_SERVER_SERVICE, NewNormalScript())
    case "uninstall":
        sc.Uninstall()
    }
    if err != nil {
        logging.Error(err)
        os.Exit(-1)
    }
    os.Exit(0)

}

func (sc *Service) SetupFor(servicePath, serviceName string, script *Script) error {

    logging.Debugf("status=begin, servicePath=%s", servicePath)

    err := utils.CreateExecutableFile(fmt.Sprintf(SERVICE_TEMPLATE, script.Script), servicePath)

    if err != nil {
        err := fmt.Sprintf("status=service-template, msg=%v", err)
        logging.Warning(err)
        return errors.New(fmt.Sprintf("status=service-template, msg=%v", err))
    }

    if utils.Exists("update-rc.d") { // debian
        _, err, _ = utils.Exec("update-rc.d", serviceName, "defaults")
        if err != nil {
            logging.Errorf("status=fatal-install-service, service=update-rc.d, msg=%s", err.Error(), err)
            os.Exit(-1)
        }
    } else if utils.Exists("chkconfig") { // redhat
        _, err, _ = utils.Exec("chkconfig", serviceName, "on")
        if err != nil {
            logging.Errorf("status=fatal-install-service, service=chkconfig, msg=%s", err.Error(), err)
            os.Exit(-2)
        }
    } else { // not known
        logging.Warningf("status=impossible to setup to start at boot")
    }

    out, err, _ := utils.Exec("service", serviceName, "stop")
    if err != nil {
        logging.Debugf("status=stop-service, msg=out=%s", string(out))
    }
    _, err, _ = utils.Exec("service", serviceName, "start")
    if err != nil {
        err := fmt.Sprintf("status=start-service, msg=%v", err)
        logging.Warning(err)
        return errors.New(err)
    }
    logging.Infof("status=success, servicePath=%s", servicePath)
    return nil

}


func (sc *Service) Uninstall() error {

    logging.Infof("status=begin")
    var err error

    if out, err, _ := utils.Exec("service", DNS_PROXY_SERVER_SERVICE, "stop"); err != nil {
        logging.Infof("status=stop-fail, msg=maibe-no-running, out=%s", string(out))
    }

    if utils.Exists("update-rc.d") {
        _, err, _ = utils.Exec("update-rc.d", "-f", DNS_PROXY_SERVER_SERVICE, "remove")
    } else if utils.Exists("chkconfig") {
        _, err, _ = utils.Exec("chkconfig", DNS_PROXY_SERVER_SERVICE, "off")
    } else {
        logging.Warningf("status=impossible to remove service")
    }
    if err != nil {
        err := fmt.Sprintf("status=fatal-remove-service, msg=%v", err)
        logging.Warning(err)
        return errors.New(err)
    }
    logging.Infof("status=success")
    return nil
}



const SERVICE_TEMPLATE = `#!/bin/sh
### BEGIN INIT INFO
# Provides:          dns-proxy-server
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       DNS PROXY SERVER
### END INIT INFO

SCRIPT=%s
RUNAS=root

PIDFILE=/var/run/dns-proxy-server.pid
LOGFILE=/var/log/dns-proxy-server.log

start() {
  if [ -f /var/run/$PIDNAME ] && kill -0 $(cat /var/run/$PIDNAME); then
    echo 'Service already running' >&2
    return 1
  fi
  echo 'Starting service…' >&2
  local CMD="$SCRIPT &> \"$LOGFILE\" & echo \$!"
  su -c "$CMD" $RUNAS > "$PIDFILE"
  echo 'Service started' >&2
}

stop() {
  if [ ! -f "$PIDFILE" ] || ! kill -0 $(cat "$PIDFILE"); then
    echo 'Service not running' >&2
    return 1
  fi
  echo 'Stopping service…' >&2
  kill -15 $(cat "$PIDFILE") && rm -f "$PIDFILE"
  echo 'Service stopped' >&2
}

uninstall() {
  echo -n "Are you really sure you want to uninstall this service? That cannot be undone. [yes|No] "
  local SURE
  read SURE
  if [ "$SURE" = "yes" ]; then
    stop
    rm -f "$PIDFILE"
    echo "Notice: log file is not be removed: '$LOGFILE'" >&2
    update-rc.d -f dns-proxy-server remove
    rm -fv "$0"
  fi
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  uninstall)
    uninstall
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|uninstall}"
esac
`
