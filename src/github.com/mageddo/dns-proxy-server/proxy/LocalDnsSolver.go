package proxy

import (
	"github.com/miekg/dns"
	"errors"
)

type LocalDnsSolver struct {

}

func (*LocalDnsSolver) Solve(name string) (*dns.Msg, error) {
	// procura no json local ou base sqlite
	return nil, errors.New("not implemented")
}
