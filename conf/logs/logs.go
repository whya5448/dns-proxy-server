package logs

import (
	"os"
	"strings"
	"github.com/mageddo/dns-proxy-server/utils/env"
)

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



func LogFile() string {
	f := GetString(os.Getenv(env.MG_LOG_FILE), flags.LogFile())
	if strings.ToLower(f) == "true" {
		return "/var/log/dns-proxy-server.log"
	}
	if strings.ToLower(f) == "false" {
		return ""
	}
	return f
}
