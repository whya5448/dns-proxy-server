package controller

import (
	"testing"
	"github.com/go-resty/resty"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/go-logging"
	"github.com/mageddo/dns-proxy-server/utils"
)

func TestGetHostnames(t *testing.T) {

	defer local.ResetConf()

	ctx := logging.NewContext()
	local.LoadConfiguration(ctx)

	err := utils.WriteToFile(`{ "remoteDnsServers": [], "envs": [
		{ "name": "MyEnv", "hostnames": [{"hostname": "github.io", "ip": [1,2,3,4], "ttl": 55}] }
	]}`, utils.GetPath(*flags.ConfPath))

	s := httptest.NewServer(nil)
	defer s.Close()

	r, err := resty.R().
		SetQueryParam("env", "MyEnv").
		Get(s.URL + HOSTNAME)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode())
	assert.Equal(t, `{"name":"MyEnv","hostnames":[{"id":1,"hostname":"github.io","ip":[1,2,3,4],"ttl":55,"env":""}]}`, r.String())

}
