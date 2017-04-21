package exitcodes

import (
	"os"
	"github.com/mageddo/log"
)

const (
	SUCCESS = iota
	FAIL_SET_DNS_AS_DEFAULT = iota
	FAIL_START_WEB_SERVER = iota
	FAIL_START_DNS_SERVER = iota
)

func Exit(code int){
	log.Logger.Errorf("m=Exit, status=exiting, code=%d", code)
	//intCode, _ := code.(int)
	os.Exit(code)
}
