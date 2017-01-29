package proxy

import (
	"errors"
	"github.com/miekg/dns"
)

type LocalDnsSolver struct {

}

func (LocalDnsSolver) Solve(question dns.Question) (*dns.Msg, error) {
	// procura no json local ou base sqlite
	return nil, errors.New("not implemented")
}
