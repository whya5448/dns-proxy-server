package exitcodes

import (
	"os"
	. "github.com/mageddo/dns-proxy-server/log"
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
	LOGGER.Errorf("m=Exit, status=exiting, code=%d", code)
	if code != SUCCESS {
		utils.Sig <- syscall.Signal(code)
		LOGGER.Info("m=Exit, status=msg-posted")
	} else {
		os.Exit(code)
	}
}
