package conf

import (
	"fmt"
	"github.com/mageddo/log"
)

type setMachineDNSServerHandler struct {
	serverIP string
}

func (hd *setMachineDNSServerHandler) process(line string, entryType DnsEntry) *string {

	switch entryType {
	case PROXY:
		log.Logger.Infof("m=SetMachineDNSServer, status=found-dns-proxy-entry")
		v := getDNSLine(hd.serverIP)
		return &v
	case SERVER:
		v := fmt.Sprintf("# %s", line)
		return &v
	default:
		return &line
	}
}

func (hd *setMachineDNSServerHandler) afterProcess(hasContent bool, foundDnsProxy bool) *string {
	if !hasContent || !foundDnsProxy {
		v := getDNSLine(hd.serverIP)
		return &v
	}
	return nil
}

func newSetMachineDnsServerHandler(serverIP string) DnsHandler {
	hd := setMachineDNSServerHandler{}
	hd.serverIP = serverIP;
	return &hd
}
