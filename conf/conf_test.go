package conf

import (
	"testing"
	"fmt"
	"os"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/go-logging"
	"github.com/mageddo/dns-proxy-server/utils"
	"flag"
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

	os.OpenFile(TMP_RESOLV_FILE, os.O_TRUNC | os.O_CREATE, 0666)
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)
	err := SetMachineDNSServer("9.9.9.9")
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile(TMP_RESOLV_FILE)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "nameserver 9.9.9.9 # dps-entry\n", string(bytes))

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

	assert.Equal(t, `# Provided by test
# nameserver 8.8.8.8 # dps-comment
nameserver 9.9.9.9 # dps-entry
`, string(bytes))

}
func TestSetMachineDNSServer_WithPreviousDnsProxyServerAndCommentSuccess(t *testing.T) {

	const TMP_RESOLV_FILE = "/tmp/test-resolv.conf"
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)

	err := ioutil.WriteFile(TMP_RESOLV_FILE, []byte("# Provided by test\n# nameserver 7.7.7.7\nnameserver 8.8.8.8\n# nameserver 10.10.10 # dps-entry"), 0666)
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

	assert.Equal(t, `# Provided by test
# nameserver 7.7.7.7
# nameserver 8.8.8.8 # dps-comment
nameserver 9.9.9.9 # dps-entry
`, string(bytes))

}


func TestRestoreResolvconfToDefault_Success(t *testing.T) {
	const TMP_RESOLV_FILE = "/tmp/test-resolv.conf"
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)

	err := ioutil.WriteFile(TMP_RESOLV_FILE, []byte("# Provided by test\n# nameserver 7.7.7.7\n# nameserver 8.8.8.8 # dps-comment\nnameserver 9.9.9.9 # dps-entry"), 0666)
	if err != nil {
		t.Error(err)
	}

	err = RestoreResolvconfToDefault()
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile(TMP_RESOLV_FILE)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

	assert.Equal(t, `# Provided by test
# nameserver 7.7.7.7
nameserver 8.8.8.8
`, string(bytes))
}

func TestRestoreResolvconfToDefault_ConfFileAlreadyOk(t *testing.T) {
	const TMP_RESOLV_FILE = "/tmp/test-resolv.conf"
	os.Setenv(env.MG_RESOLVCONF, TMP_RESOLV_FILE)

	originalFileContent := "# Provided by test\n# nameserver 8.8.8.8\nnameserver 9.9.9.9\n"
	err := ioutil.WriteFile(TMP_RESOLV_FILE, []byte(originalFileContent), 0666)
	if err != nil {
		t.Error(err)
	}

	err = RestoreResolvconfToDefault()
	if err != nil {
		t.Error(err)
	}
	bytes, err := ioutil.ReadFile(TMP_RESOLV_FILE)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bytes))

	assert.Equal(t, originalFileContent, string(bytes))
}


func TestDefaultFlagValues(t *testing.T) {
	assert.Equal(t, *flags.WebServerPort, WebServerPort())
	assert.Equal(t, *flags.DnsServerPort, DnsServerPort())
	assert.Equal(t, *flags.SetupResolvconf, SetupResolvConf())
}

func TestFlagValuesFromArgs(t *testing.T) {
	os.Args = []string{"cmd", "-web-server-port=8282", "--server-port=61", "-default-dns=false"}
	flag.Parse()
	assert.Equal(t, 8282, WebServerPort())
	assert.Equal(t, 61, DnsServerPort())
	assert.Equal(t, false, SetupResolvConf())
}

func TestFlagValuesFromConf(t *testing.T) {

	defer local.ResetConf()
	ctx := logging.NewContext()
	local.LoadConfiguration(ctx)

	err := utils.WriteToFile(`{ "webServerPort": 8080, "dnsServerPort": 62, "defaultDns": false }`, utils.GetPath(*flags.ConfPath))
	assert.Nil(t, err)

	assert.Equal(t, 8080, WebServerPort())
	assert.Equal(t, 62, DnsServerPort())
	assert.Equal(t, false, SetupResolvConf())
}
