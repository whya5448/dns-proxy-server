package utils

import (
	. "github.com/mageddo/dns-proxy-server/log"
	"os/exec"
	"syscall"
	"reflect"
)

func Exec(cmd string, args ...string) ( out []byte, err error, exitCode int ){

	LOGGER.Infof("cmd=%s, args=%v", cmd, args)

	execution := exec.Command(cmd, args...)
	// ja chama o run dentro dele
	out, err = execution.CombinedOutput()

	if err != nil {
		LOGGER.Infof("status=error, type=%v, err=%v", reflect.TypeOf(err), err)
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
			return
		} else {
			exitCode = -255
			return
		}
	} else {
		exitCode = execution.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}

	if exitCode != 0 {
		LOGGER.Warningf("status=bad-exit-code, status=%d", exitCode)
		return
	}
	LOGGER.Infof("status=success, cmd=%s", cmd)
	return
}

func Exists(cmd string) bool {
	_, _, i := Exec("sh", "-c", "command -v " + cmd + " || false")
	switch i {
	case 0:
		return true
	case -255:
		panic("Command verificator not exists")
	default:
		return false

	}
}
