package proxy

import (
	"testing"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/mock"
	"github.com/mageddo/dns-proxy-server/cache/lru"
	hashlru "github.com/hashicorp/golang-lru"
	"fmt"
	"github.com/mageddo/dns-proxy-server/cache/store"
)

type MockCache struct {
	mock.Mock
	lru.LRUCache
}

func (m *MockCache) Get(key interface{}) interface{} {
	fmt.Println("mocked!!!!!")
	m.Called(key)
	return m.LRUCache.Get(key)
	//return args.Get(0)
}

func TestLocalDnsSolver_Solve(t *testing.T) {

	ctx := logging.NewContext()
	conf, err := local.LoadConfiguration(ctx)
	assert.Nil(t, err, "failed to load configuration")

	expectedHostname := "github.com"
	host := local.HostnameVo{Hostname: expectedHostname, Env:"", Ttl:50, Ip:[4]byte{192,168,0,1}}
	conf.AddHostname(ctx, "", host)

	question := new(dns.Question)
	question.Name = expectedHostname + "."
	solver := NewLocalDNSSolver(store.GetInstance())

	// act
	res, err := solver.Solve(ctx, *question)
	assert.Nil(t, err, "Fail to solve")

	// assert
	assert.Equal(t, 1, len(res.Answer))
	assert.Equal(t, "github.com.	0	IN	A	192.168.0.1", res.Answer[0].String())

}

func TestLocalDnsSolver_SolveValidatingCache(t *testing.T) {

	ctx := logging.NewContext()
	conf, err := local.LoadConfiguration(ctx)
	assert.Nil(t, err, "failed to load configuration")

	expectedHostname := "github.com"
	host := local.HostnameVo{Hostname: expectedHostname, Env:"", Ttl:50, Ip:[4]byte{192,168,0,1}}
	conf.AddHostname(ctx, "", host)

	question := new(dns.Question)
	question.Name = expectedHostname + "."

	mockCache := &MockCache{}
	mockCache.Cache, err = hashlru.New(1)
	assert.Nil(t, err, "Failed to create cache")
	solver := NewLocalDNSSolver(mockCache)

	mockCache.On("Get", expectedHostname).Twice()


	for i:=2; i > 0; i-- {

		// act
		res, err := solver.Solve(ctx, *question)
		assert.Nil(t, err, "Fail to solve")

		// assert
		assert.Equal(t, 1, len(res.Answer))
		assert.Equal(t, "github.com.	0	IN	A	192.168.0.1", res.Answer[0].String())

	}

}
