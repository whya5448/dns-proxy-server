package proxy

import "github.com/miekg/dns"

type DnsSolver interface {

	Solve(question dns.Question) (*dns.Msg, error)


}
