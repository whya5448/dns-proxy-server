package exitcodes

import (
	"os"
	"github.com/mageddo/log"
	"github.com/mageddo/dns-proxy-server/utils"
	"syscall"
)

const (
	SUCCESS = iota
	FAIL_SET_DNS_AS_DEFAULT = iota
	FAIL_START_WEB_SERVER = iota
	FAIL_START_DNS_SERVER = iota
)

func Exit(code int){
	log.Logger.Errorf("m=Exit, status=exiting, code=%d", code)
	if code != SUCCESS {
		utils.Sig <- syscall.Signal(code)
	} else {
		os.Exit(code)
	}
}
