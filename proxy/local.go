package proxy

import (
	. "github.com/mageddo/dns-proxy-server/log"
	"errors"
	"github.com/miekg/dns"
	"net"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/cache"
)

type LocalDnsSolver struct {
	Cache cache.Cache
}

func (s LocalDnsSolver) Solve(ctx context.Context, question dns.Question) (*dns.Msg, error) {

	key := question.Name[:len(question.Name)-1]
	var hostname *local.HostnameVo
	if s.Cache.ContainsKey(key) {
		LOGGER.Debugf("status=from-cache, key=%s, value=%v", key, s.Cache.Get(key))
		if s.Cache.Get(key) != nil {
			hostname = s.Cache.Get(key).(*local.HostnameVo)
		}
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
		val := s.Cache.PutIfAbsent(key, hostname);
		LOGGER.Debugf("status=put, key=%s, value=%v", key, val)
	}

	if hostname != nil {
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A: net.IPv4(hostname.Ip[0], hostname.Ip[1], hostname.Ip[2], hostname.Ip[3]),
		}

		m := new(dns.Msg)
		m.Answer = append(m.Answer, rr)
		LOGGER.Debugf("status=success, solver=local, key=%s", key)
		return m, nil
	}
	return nil, errors.New("hostname not found " + key)
}

func NewLocalDNSSolver(c cache.Cache) *LocalDnsSolver {
	return &LocalDnsSolver{c}
}
