package conf


type dnsServerCleanerHandler struct {
	serverIP string
}

func (hd dnsServerCleanerHandler) process(line string, entryType DnsEntry) *string {

	switch entryType {
	case PROXY:
		return nil
	case SERVER:
		panic("it can not happen")
	case COMMENTED_SERVER:
		v := line[2:]
		return &v
	default:
		return &line
	}
}

func (hd dnsServerCleanerHandler) afterProcess(hasContent bool, foundDnsProxy bool) *string {
	if !hasContent || !foundDnsProxy {
		v := getDNSLine(hd.serverIP)
		return &v
	}
	return nil
}

func newDNSServerCleanerHandler() DnsHandler {
	hd := dnsServerCleanerHandler{}
	return &hd
}
