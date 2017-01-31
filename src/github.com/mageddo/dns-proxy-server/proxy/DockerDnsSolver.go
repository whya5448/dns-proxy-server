package proxy

import (
	"github.com/miekg/dns"
	"github.com/mageddo/dns-proxy-server/events"
	"github.com/mageddo/log"
	"net"
	"strings"
	"strconv"
	"errors"
)

type DockerDnsSolver struct {

}

func (DockerDnsSolver) Solve(question dns.Question) (*dns.Msg, error) {

	log.Logger.Infof("m=solve, status=begin, solver=docker, name=%s", question.Name)

	if events.ContainsKey(question.Name) {

		ip := events.Get(question.Name)
		ipArr := strings.Split(ip, ".")
		i1, _ := strconv.Atoi(ipArr[0])
		i2, _ := strconv.Atoi(ipArr[1])
		i3, _ := strconv.Atoi(ipArr[2])
		i4, _ := strconv.Atoi(ipArr[3])

		m := new(dns.Msg)
		m.Answer = dns.A{
			Hdr: dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A: net.IPv4(i1, i2, i3, i4),
		}
		return m, nil
	}
	return nil, errors.New("hostname not found")

	log.Logger.Infof("m=solve, status=success, solver=docker")

}
