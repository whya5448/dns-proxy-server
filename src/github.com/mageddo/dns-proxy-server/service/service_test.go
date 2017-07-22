package service

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/log"
	"fmt"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/dns-proxy-server/utils"
	"strings"
)

func TestSetupFor_NormalModeInstallStartSuccess(t *testing.T) {

	if !flags.IsTestVersion() {
		log.Logger.Infof("status=test-skiped")
		return
	}

	ctx := log.GetContext()

	sc := NewService(ctx)
	cmd := "'bash -c \"echo hi && sleep 20\"'"
	err := sc.SetupFor(DNS_PROXY_SERVER_PATH, DNS_PROXY_SERVER_SERVICE, &Script{cmd})
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadFile(DNS_PROXY_SERVER_PATH)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf(SERVICE_TEMPLATE, cmd), string(bytes))

	out, err, code := utils.Exec("bash", "-c", "ps aux | grep \"echo hi\" && exit ")
	assert.Equal(t, 0, code)
	assert.Nil(t, err)

	str := string(out)
	assert.NotEqual(t, -1 , strings.Index(str, "echo hi && sleep 20"))

}
