package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
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