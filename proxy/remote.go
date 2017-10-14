package proxy

import (
	"github.com/miekg/dns"
	"net"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/events/local"
	log "github.com/mageddo/go-logging"

	"github.com/mageddo/dns-proxy-server/cache/store"
)

const SERVERS = "SERVERS"

type RemoteDnsSolver struct {
	confloader func(ctx context.Context) (*local.LocalConfiguration, error)
}

// reference https://miek.nl/2014/August/16/go-dns-package/
func (r RemoteDnsSolver) Solve(ctx context.Context, question dns.Question) (*dns.Msg, error) {
	c := store.GetInstance()
	logger := log.NewLog(ctx)
	client := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(question.Name), question.Qtype) // CAN BE A, AAA, MX, etc.
	m.RecursionDesired = true

	var config *local.LocalConfiguration
	var err error
	if !c.ContainsKey(SERVERS) {
		logger.Debugf("status=hot-load")
		if config, err = r.confloader(ctx); err != nil {
			logger.Errorf("error=%v",err)
			return nil, err
		}
		c.PutIfAbsent(SERVERS, config)
	} else {
		logger.Debugf("status=from-cache")
	}
	config = c.Get(SERVERS).(*local.LocalConfiguration)

	for _, server := range config.GetRemoteServers(ctx) {

		if len(server) != 4 {
			logger.Warning("status=wrong-server, server=%+v", server)
			continue
		}

		// server and port to ask
		formatServer := fmt.Sprintf("%d.%d.%d.%d", server[0], server[1], server[2], server[3])
		logger.Debugf("status=format-server, server=%s", formatServer)

		var r *dns.Msg
		r, _, err = client.Exchange(m, net.JoinHostPort(formatServer, "53"))

		// if the answer not be returned
		if r == nil {
			err = errors.New(fmt.Sprintf("status=answer-can-not-be-null, err=%v", err))
			logger.Infof("status=no-answer, err=%s", err)
			continue
		} else if r.Rcode != dns.RcodeSuccess { // what the code of the return message ?
			err = errors.New(fmt.Sprintf("status=invalid-answer-name, name=%s, rcode=%d", question.Name, r.Rcode))
			logger.Infof("status=bad-code, name=%s, rcode=%d, err=%s", question.Name, r.Rcode, err)
			continue
		}
		return r, nil
	}
	return nil, err
}

func NewRemoteDnsSolver() *RemoteDnsSolver {
	return &RemoteDnsSolver{
		confloader: func(ctx context.Context) (*local.LocalConfiguration, error) {
		return local.LoadConfiguration(ctx)
	}}
}
