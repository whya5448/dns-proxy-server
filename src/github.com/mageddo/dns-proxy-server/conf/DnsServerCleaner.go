package conf


type dnsServerCleanerHandler struct {
	serverIP string
}

func (hd dnsServerCleanerHandler) process(line string, entryType DnsEntry) *string {

	switch entryType {
	case PROXY:
		return nil
	case COMMENTED_SERVER:
		v := line[2:]
		return &v
	case SERVER:
		panic("it can not happen")
	default:
		return &line
	}
}

func (hd dnsServerCleanerHandler) afterProcess(hasContent bool, foundDnsProxy bool) *string {
	return nil
}

func newDNSServerCleanerHandler() DnsHandler {
	hd := dnsServerCleanerHandler{}
	return &hd
}
