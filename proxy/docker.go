package proxy

import (
	"errors"
	"github.com/mageddo/dns-proxy-server/cache"
	log "github.com/mageddo/go-logging"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
	"net"
	"strconv"
	"strings"
)

type DockerDnsSolver struct {
	c cache.Cache
}

func (s DockerDnsSolver) Solve(ctx context.Context, question dns.Question) (*dns.Msg, error) {

	logger := log.NewLog(ctx)
	key := question.Name[:len(question.Name)-1]
	if s.c.ContainsKey(key) {
		logger.Debugf("solver=docker, status=solved-key, solver=docker, hostname=%s, ip=%+v", key, s.c.Get(key))
		return s.getMsg(key, question), nil
	}
	i := strings.Index(key, ".")
	if i > 0 && s.c.ContainsKey(key[i:]) {
		logger.Debugf("solver=docker, status=solved-key-wildcard, solver=docker, hostname=%s, ip=%+v", key, s.c.Get(key[i:]))
		return s.getMsg(key[i:], question), nil
	}
	return nil, errors.New("hostname not found " + key)
}

func NewDockerSolver(c cache.Cache) DockerDnsSolver {
	return DockerDnsSolver{c}
}

func (s DockerDnsSolver) getMsg(key string, question dns.Question) *dns.Msg {
	ip := s.c.Get(key).(string)
	ipArr := strings.Split(ip, ".")
	i1, _ := strconv.Atoi(ipArr[0])
	i2, _ := strconv.Atoi(ipArr[1])
	i3, _ := strconv.Atoi(ipArr[2])
	i4, _ := strconv.Atoi(ipArr[3])

	rr := &dns.A{
		Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.IPv4(byte(i1), byte(i2), byte(i3), byte(i4)),
	}

	m := new(dns.Msg)
	m.Answer = append(m.Answer, rr)
	return m
}
