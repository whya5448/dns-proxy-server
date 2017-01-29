package dns

import "github.com/miekg/dns"

type LocalDnsSolver struct {

}

func (*LocalDnsSolver) Solve(name string) *dns.Msg {
	// procura no json local ou base sqlite
	return nil
}
