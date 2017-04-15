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

func TestSetMachineDNSServer_EmptyFileSuccess(t *testing.T) {

	const TMP_RESOLV_FILE = "/tmp/test-resolv.conf"

	os.OpenFile(TMP_RESOLV_FILE, os.O_TRUNC, 0666)
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)
	err := SetMachineDNSServer("9.9.9.9")
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile(TMP_RESOLV_FILE)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

	// Output: nameserver 9.9.9.9 # dns-proxy-server

}

func TestSetMachineDNSServer_WithPreviousDnsServerAndCommentSuccess(t *testing.T) {

	const TMP_RESOLV_FILE = "/tmp/test-resolv.conf"
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)

	err := ioutil.WriteFile(TMP_RESOLV_FILE, []byte("# Provided by test\nnameserver 8.8.8.8"), 0666)
	if err != nil {
		t.Error(err)
	}

	err = SetMachineDNSServer("9.9.9.9")
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile(TMP_RESOLV_FILE)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

	AssertEqual(t, `nameserver 9.9.9.9 # dns-proxy-server
	# Provided by test
	nameserver 8.8.8.8`, string(bytes), "not match")

}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}