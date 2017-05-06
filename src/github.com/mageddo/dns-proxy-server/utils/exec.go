package utils

import (
	"os/exec"
	"errors"
	"fmt"
	"syscall"
	"github.com/mageddo/log"
)

func Exec(cmd string, args ...string) ( out []byte, err error, exitCode int ){

	log.Logger.Infof("m=Exec, cmd=%s, args=%v", cmd, args)

	execution := exec.Command(cmd, args...)

	// ja chama o run dentro dele
	out, err = execution.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
			err = errors.New(fmt.Sprintf("m=Exec, exitcode=%d, err=%s, out=%s", exitCode, err.Error(), string(out)))
			return
		}
	} else {
		exitCode = execution.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}

	if exitCode != 0 {
		log.Logger.Warningf("m=Exec, status=bad-exit-code, status=%d", exitCode)
		return
	}
	log.Logger.Infof("m=Exec, status=success, cmd=%s", cmd)
	return
}

func Exists(cmd string) bool {
	_, _, i := Exec("hash", cmd, "||", "false")
	return i == 0
}
