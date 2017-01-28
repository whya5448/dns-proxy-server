package utils

import (
	"github.com/miekg/dns"
	"net"
	"github.com/mageddo/log"
	
)

// reference https://miek.nl/2014/August/16/go-dns-package/
func SolveName(hostname string) *dns.Msg {


	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeA) // CAN BE A, AAA, MX, etc.
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port)) // server and port to ask

	// if the answer not be returned
	if r == nil {
		log.Logger.Fatalf("**** error: %s", err.Error())
	}

	// what the code of the return message ?
	if r.Rcode != dns.RcodeSuccess {
		log.Logger.Fatalf(" *** invalid answer name %s after MX query for %s", hostname, hostname)
	}

	// looping through the anwsers
	return r

}
