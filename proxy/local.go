package proxy

import (
	. "github.com/mageddo/dns-proxy-server/log"
	"errors"
	"github.com/miekg/dns"
	"net"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/cache/lru"
	c "github.com/mageddo/dns-proxy-server/cache"
)

type LocalDnsSolver struct {}

var cache c.Cache;
func init(){
	cache = lru.New(256);
}

func (LocalDnsSolver) Solve(ctx context.Context, question dns.Question) (*dns.Msg, error) {

	key := question.Name[:len(question.Name)-1]
	var hostname *local.HostnameVo
	if cache.ContainsKey(key) {
		LOGGER.Debugf("status=from-cache, key=%s, value=%v", key, hostname.Hostname)
		hostname = cache.Get(key).(*local.HostnameVo)
	} else {
		LOGGER.Debugf("status=hot-load, key=%s", key)
		conf, err := local.LoadConfiguration(ctx)
		if err != nil {
			LOGGER.Errorf("status=could-not-load-conf, err=%v", err)
			return nil, err
		}
		activeEnv,_ := conf.GetActiveEnv()
		if activeEnv == nil {
			return nil, errors.New("original env")
		}
		hostname,_ = activeEnv.GetHostname(key)
		val := cache.PutIfAbsent(key, hostname);
		LOGGER.Debugf("status=put, key=%s, value=%v", key, val)
	}

	if hostname != nil {
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A: net.IPv4(hostname.Ip[0], hostname.Ip[1], hostname.Ip[2], hostname.Ip[3]),
		}

		m := new(dns.Msg)
		m.Answer = append(m.Answer, rr)
		LOGGER.Debugf("status=success, solver=local")
		return m, nil
	}
	return nil, errors.New("hostname not found")
}
