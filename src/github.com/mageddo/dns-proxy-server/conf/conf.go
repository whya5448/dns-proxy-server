package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"bytes"
	"os"
	"bufio"
	"strings"
	"net"
	"github.com/mageddo/log"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"io/ioutil"
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

func GetString(value, defaultValue string) string {

	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func SetMachineDNSServer(serverIP string) error {

	var newResolvConfBuff bytes.Buffer

	log.Logger.Infof("m=SetMachineDNSServer, status=begin")

	resolvconf := GetString(os.Getenv(env.MG_RESOLVCONF), "/etc/resolv.conf")
	fileRead, err := os.Open(resolvconf)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	log.Logger.Infof("m=SetMachineDNSServer, status=open-conf-file, file=%s", fileRead.Name())

	scanner := bufio.NewScanner(fileRead)
	var (
		hasContent = false
		foundDnsProxyEntry = false
	)

	for scanner.Scan() {
		hasContent = true
		line := scanner.Text()
		if strings.HasSuffix(line, "# dns-proxy-server") {

			// this line is dns proxy server nameserver entry
			log.Logger.Infof("m=SetMachineDNSServer, status=found-dns-proxy-entry")
			newResolvConfBuff.WriteString(getDNSLine(serverIP))
			foundDnsProxyEntry = true

		} else if strings.HasPrefix(line, "#") {

			// linha comentada
			log.Logger.Infof("m=SetMachineDNSServer, status=commented-line")
			newResolvConfBuff.WriteString(line)

		} else if strings.HasPrefix(line, "nameserver") {

			log.Logger.Infof("m=SetMachineDNSServer, status=nameserver-line")
			newResolvConfBuff.WriteString("# " + line)

		} else {

			log.Logger.Infof("m=SetMachineDNSServer, status=else-line")
			newResolvConfBuff.WriteString(line)

		}
		newResolvConfBuff.WriteByte('\n')
	}
	if !hasContent || !foundDnsProxyEntry {
		newResolvConfBuff.WriteString(getDNSLine(serverIP))
	}
	stats, _ := fileRead.Stat()
	length := newResolvConfBuff.Len()
	err = ioutil.WriteFile(resolvconf, newResolvConfBuff.Bytes(), stats.Mode())
	if err != nil {
		return err
	}
	log.Logger.Infof("m=SetMachineDNSServer, status=success, buffLength=%d", length)
	return nil
}
func getDNSLine(serverIP string) string {
	return "nameserver " + serverIP + " # dns-proxy-server"
}

func SetCurrentDNSServerToMachine() error {

	log.Logger.Infof("m=SetCurrentDNSServerToMachine, status=begin")
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