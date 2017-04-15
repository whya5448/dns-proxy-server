package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"bytes"
	"fmt"
	"os"
	"bufio"
	"strings"
	"net"
)

func CpuProfile() string {
return *flags.Cpuprofile
}

func Compress() bool {
return *flags.Compress
}

func Tsig() string {
return *flags.Tsig
}

func WebServerPort() int {
port := local.GetConfigurationNoCtx().WebServerPort
if port <= 0 {
return *flags.WebServerPort
}
return port
}

func DnsServerPort() int {
	port := local.GetConfigurationNoCtx().DnsServerPort
	if port <= 0 {
	return *flags.DnsServerPort
	}
	return port
}

func SetupResolvConf() bool {
	return *flags.SetupResolvconf
}

func ConfPath() string {
	return *flags.ConfPath
}

func SetMachineDNSServer(serverIP string) error {

	var newResolvConfBuff bytes.Buffer

	fmt.Println(newResolvConfBuff.String())
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasSuffix(line, "# dns-proxy-server") { // this line is dns proxy server nameserver entry
			newResolvConfBuff.WriteString(serverIP + "# dns-proxy-server")
		}else if strings.HasPrefix(line, "#") { // linha comentada
			newResolvConfBuff.WriteString(line)
		} else if strings.HasPrefix(line, "nameserver") {
			newResolvConfBuff.WriteString("# " + line)
		} else {
			newResolvConfBuff.WriteString(line)
		}
		newResolvConfBuff.WriteByte('\n')

	}
	newResolvConfBuff.WriteTo(bufio.NewWriter(file))
	return nil
}

func SetCurrentDNSServerToMachine(serverIP string) error {

	ip, err := getCurrentIpAddress()
	if err != nil {
		return err
	}
	return SetMachineDNSServer(ip)
}

func getCurrentIpAddress() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ip := addr.String()
		if strings.Contains(ip, "/") {
			if !strings.HasPrefix(ip, "127") {
				return ip[:strings.Index(ip, "/")], nil
			}
		}
	}
	return "", nil

}