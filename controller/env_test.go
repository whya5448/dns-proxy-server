package controller

import (
	"testing"
	"net/http/httptest"
	"github.com/go-resty/resty"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/go-logging"
)

func TestGetActiveEnvSuccess(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	r, err := resty.R().Get(s.URL + ENV_ACTIVE)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, "{\n\t\"name\": \"\"\n}", r.String())

}

func TestPutChangeActiveEnvThatDoesNotExistsError(t *testing.T) {

	s := httptest.NewServer(nil)
	defer s.Close()

	r, err := resty.R().
		SetBody(`{"name": "testEnv"}`).
		Put(s.URL + ENV_ACTIVE)

	assert.Nil(t, err)
	assert.Equal(t, 400, r.StatusCode())
	assert.Equal(t, `{"code":400,"message":"Env not found: testEnv"}`, r.String())

}

func TestPutChangeActiveEnvSuccess(t *testing.T) {

	defer local.ResetConf()

	ctx := logging.NewContext()
	local.LoadConfiguration(ctx)

	err := utils.WriteToFile(`{
	"remoteDnsServers": [], "envs": [
		{ "name": "testEnv",
		"hostnames": [
			{
				"id": 1,
				"hostname": "mageddo.com",
				"ip": [192, 168, 0, 3],
				"ttl": 255
			}]
		}
	]}`, utils.GetPath(*flags.ConfPath))

	assert.Nil(t, err)

	s := httptest.NewServer(nil)
	defer s.Close()

	r, err := resty.R().
		SetBody(`{"name": "testEnv"}`).
		Put(s.URL + ENV_ACTIVE)

	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Empty(t, r.String())

	r, err = resty.R().Get(s.URL + ENV_ACTIVE)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode())
	assert.Equal(t, "{\n\t\"name\": \"testEnv\"\n}", r.String())

}
