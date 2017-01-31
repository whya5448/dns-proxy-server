package proxy

import (
	"github.com/miekg/dns"
	"net"
	"errors"
	"fmt"
)

type RemoteDnsSolver struct {

}

// reference https://miek.nl/2014/August/16/go-dns-package/
func (RemoteDnsSolver) Solve(question dns.Question) (*dns.Msg, error) {

		//config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
		c := new(dns.Client)

		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(question.Name), question.Qtype) // CAN BE A, AAA, MX, etc.
		m.RecursionDesired = true

		//r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port)) // server and port to ask
		r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53")) // server and port to ask

		// if the answer not be returned
		if r == nil {
			return nil, errors.New(fmt.Sprintf("status=answer-can-not-be-bull, err=%v", err))
		}

		// what the code of the return message ?
		if r.Rcode != dns.RcodeSuccess {
			return nil, errors.New(fmt.Sprintf("status=invalid-answer-name, name=%s, rcode=%d", question.Name, r.Rcode))
		}

		return r, nil

	}
