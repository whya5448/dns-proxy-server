package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"os"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"github.com/mageddo/go-logging"
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
	if conf, _ := getConf(); conf != nil && conf.WebServerPort > 0 {
		return conf.WebServerPort
	}
	return *flags.WebServerPort
}

func DnsServerPort() int {
	if conf, _ := getConf(); conf != nil && conf.DnsServerPort > 0 {
		return conf.DnsServerPort
	}
	return *flags.DnsServerPort
}

func SetupResolvConf() bool {
	if conf, _ := getConf(); conf != nil && conf.DefaultDns != nil {
		return *conf.DefaultDns
	}
	return *flags.SetupResolvconf
}

func GetResolvConf() string {
	return GetString(os.Getenv(env.MG_RESOLVCONF), "/etc/resolv.conf")
}

func getConf() (*local.LocalConfiguration, error) {
	return local.LoadConfiguration(logging.NewContext())
}

func GetString(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func LogLevel() string {
	if lvl := os.Getenv(env.MG_LOG_LEVEL); lvl != "" {
		return lvl
	}

	if conf, _ := getConf(); conf != nil && conf.LogLevel != "" {
		return conf.LogLevel
	}
	return flags.LogLevel()
}
