package service

import (
	"fmt"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
	"os"
	"golang.org/x/net/context"
	"errors"
)

type Service struct {
	ctx context.Context
	logger *log.IdLogger
}

type Script struct {
	Script string
}

func NewService(ctx context.Context) *Service {
	return &Service{ctx, log.GetLogger(ctx)}
}

func (sc *Service) Install() {

	setupServiceFlag, servicePath := *flags.SetupService, "/etc/init.d/dns-proxy-server"

	sc.logger.Infof("setupservice=%s, version=%s", setupServiceFlag, flags.GetRawCurrentVersion())
	var err error
	switch setupServiceFlag {
	case "docker":
		err = sc.SetupFor(servicePath, NewDockerScript())
	case "normal":
		err = sc.SetupFor(servicePath, NewNormalScript())
	case "uninstall":
		sc.Uninstall()
	}
	if err != nil {
		sc.logger.Error(err)
		os.Exit(-1)
	}
	os.Exit(0)

}

func (sc *Service) SetupFor(servicePath string, script *Script) error {

	sc.logger.Debugf("status=begin, servicePath=%s", servicePath)

	err := utils.CreateExecutableFile(SERVICE_TEMPLATE, servicePath)
	_, err, _ = utils.Exec("sed", "-i", fmt.Sprintf("s/%s/%s/g", "<SCRIPT>", script), servicePath)
	if err != nil {
		return errors.New(fmt.Sprintf("status=error-prepare-service, msg=%v", err))
	}

	if err != nil {
		err := fmt.Sprintf("status=service-template, msg=%v", err)
		sc.logger.Warning(err)
		return errors.New(fmt.Sprintf("status=service-template, msg=%v", err))
	}

	if utils.Exists("update-rc.d") { // debian
		_, err, _ = utils.Exec("update-rc.d", "dns-proxy-server", "defaults")
		if err != nil {
			sc.logger.Fatalf("status=fatal-install-service, service=update-rc.d, msg=%s", err.Error())
		}
	} else if utils.Exists("chkconfig") { // redhat
		_, err, _ = utils.Exec("chkconfig", "dns-proxy-server", "on")
		if err != nil {
			sc.logger.Fatalf("status=fatal-install-service, service=chkconfig, msg=%s", err.Error())
		}
	} else { // not known
		sc.logger.Warningf("m=ConfigSetupService, status=impossible to setup to start at boot")
	}

	out, err, _ := utils.Exec("service", "dns-proxy-server", "stop")
	if err != nil {
		sc.logger.Debugf("status=stop-service, msg=out=%s", string(out))
	}
	_, err, _ = utils.Exec("service", "dns-proxy-server", "start")
	if err != nil {
		err := fmt.Sprintf("status=start-service, msg=%v", err)
		sc.logger.Warning(err)
		return errors.New(err)
	}
	sc.logger.Infof("status=success, servicePath=%s", servicePath)
	return nil

}


func (sc *Service) Uninstall() error {

	sc.logger.Infof("status=begin")
	var err error

	if out, err, _ := utils.Exec("service", "dns-proxy-server", "stop"); err != nil {
		sc.logger.Infof("status=stop-fail, msg=maibe-no-running, out=%s", string(out))
	}

	if utils.Exists("update-rc.d") {
		_, err, _ = utils.Exec("update-rc.d", "-f", "dns-proxy-server", "remove")
	} else if utils.Exists("chkconfig") {
		_, err, _ = utils.Exec("chkconfig", "dns-proxy-server", "off")
	} else {
		sc.logger.Warningf("status=impossible to remove service")
	}
	if err != nil {
		err := fmt.Sprintf("status=fatal-remove-service, msg=%v", err)
		sc.logger.Warning(err)
		return errors.New(err)
	}
	sc.logger.Infof("status=success")
	return nil
}



const SERVICE_TEMPLATE = `
#!/bin/sh
### BEGIN INIT INFO
# Provides:          dns-proxy-server
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       DNS PROXY SERVER
### END INIT INFO

SCRIPT=<SCRIPT>
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
  retart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|uninstall}"
esac
`