package dns

import (
	"github.com/miekg/dns"
	"net"
	"errors"
	"fmt"
)

type RemoteDnsSolver struct {

}

// reference https://miek.nl/2014/August/16/go-dns-package/
func (*RemoteDnsSolver) Solve(name string) *dns.Msg {

		//config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
		c := new(dns.Client)

		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(name), dns.TypeA) // CAN BE A, AAA, MX, etc.
		m.RecursionDesired = true

		//r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port)) // server and port to ask
		r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53")) // server and port to ask

		// if the answer not be returned
		if r == nil {
		panic(err)
		}

		// what the code of the return message ?
		if r.Rcode != dns.RcodeSuccess {
		panic(errors.New(fmt.Sprintf(" *** invalid answer name %s after MX query for %s", name, name)))
		}

		// looping through the anwsers
		return r

	}
