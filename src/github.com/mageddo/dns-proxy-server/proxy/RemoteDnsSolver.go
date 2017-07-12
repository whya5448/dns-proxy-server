package proxy

import (
	"github.com/miekg/dns"
	"net"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/conf"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/log"
)

type RemoteDnsSolver struct {

}

// reference https://miek.nl/2014/August/16/go-dns-package/
func (RemoteDnsSolver) Solve(ctx context.Context, question dns.Question) (*dns.Msg, error) {

	logger := log.GetLogger(ctx)
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(question.Name), question.Qtype) // CAN BE A, AAA, MX, etc.
	m.RecursionDesired = true

	var err error
	config := local.GetConfiguration(ctx)

	for server := range config.GetRemoteServers(ctx) {

		if (len(server) != 4){
			logger.Warning("status=wrong-server, server=%+v", server)
			continue
		}

		// server and port to ask
		r, _, err := c.Exchange(m, net.JoinHostPort(fmt.Sprintf("%d.%d.%d.%d", server[0], server[1], server[2], server[3]), "53"))

			// if the answer not be returned
			if r == nil {
				err = errors.New(fmt.Sprintf("status=answer-can-not-be-bull, err=%v", err))
				logger.Infof("status=no-answer, err=%s", err)
				continue
			} else if r.Rcode != dns.RcodeSuccess { // what the code of the return message ?
				err = errors.New(fmt.Sprintf("status=invalid-answer-name, name=%s, rcode=%d", question.Name, r.Rcode))
				logger.Infof("status=bad-code, name=%d, rcode=%d, err=%s", question.Name, r.Rcode, err)
				continue
			}
			return r, nil
		}
		return nil, err
	}
