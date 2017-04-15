package conf

import (
	"testing"
	"fmt"
	"os"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"io/ioutil"
)

func TestGetCurrentIpAddress(t *testing.T){

	ip, err := getCurrentIpAddress()
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(">> " + ip)

}

func TestSetMachineDNSServer(t *testing.T) {

	//x := os.Open("/tmp/test-resolv.conf")


	os.Setenv(env.MG_RESOLVCONF, "/tmp/test-resolv.conf")
	err := SetMachineDNSServer("9.9.9.9")
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile("/tmp/test-resolv.conf")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

}