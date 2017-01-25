package main

import (
	"github.com/miekg/dns"
	"os"
	"net"
	"github.com/mageddo/log"
	
)

func main(){


	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(os.Args[1]), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))

	if r == nil {
		log.Logger.Fatalf("**** error: %s", err.Error())
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Logger.Fatalf(" *** invalid answer name %s after MX query for %s", os.Args[1], os.Args[1])
	}

	for _, a := range r.Answer {
		log.Logger.Infof("%v", a)
	}

}
