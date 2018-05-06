package local

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/dns-proxy-server/cache/store"
)

func TestSaveConfiguration_ClearCacheAfterChangeConfiguration(t *testing.T) {

	defer ResetConf()

	expectedHostname := "github.io"

	conf, err := LoadConfiguration()
	assert.Nil(t, err, "could not load conf")

	cache := store.GetInstance()

	env, _ := conf.GetActiveEnv()
	foundHostname, _ := env.GetHostname(expectedHostname)
	assert.Nil(t, foundHostname, "can't be nil")

	// setting the host
	cache.Put(expectedHostname, foundHostname)
	assert.True(t, cache.ContainsKey(expectedHostname))
	assert.Nil(t, cache.Get(expectedHostname))

	// changing value for the hostname at configuration database
	hostname := HostnameVo{Ip: [4]byte{192,168,0,2}, Ttl:30, Env:"", Hostname: expectedHostname}
	conf.AddHostname( "", hostname)

	// cache must be clear after add a hostname in conf
	assert.False(t, cache.ContainsKey(expectedHostname))
	foundHostname, _ = env.GetHostname(expectedHostname)
	assert.Equal(t, [4]byte{192,168,0,2}, foundHostname.Ip)

}
