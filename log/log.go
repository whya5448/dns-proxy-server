package log

import (
	"github.com/mageddo/go-logging"
	"os"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/dns-proxy-server/conf"
)

var LOGGER logging.Log

func init(){

	mode := "dev"
	var out = os.Stdout
	var err error

	if f := flags.LogFile(); f != "" {
		if out, err = os.OpenFile(f, os.O_CREATE|os.O_APPEND, 0766); err != nil {
			panic(err)
		}
	}
	backend := logging.NewLogBackend(out, "", 0)

	// setando o log dependendo do ambiente
	switch mode {
	case "prod":
		format := logging.MustStringFormatter(
			`%{time:06-01-02 15:04:05} %{level:.3s} %{message}`,
		)
		leveledBackend := logging.AddModuleLevel(logging.NewBackendFormatter(backend, format))
		logging.SetBackend(leveledBackend)
		logging.SetLevel(logging.INFO, "")
		break
	default:
		format := logging.MustStringFormatter(
			`%{color}%{time:06-01-02 15:04:05.000} %{level:.3s} %{color:reset}%{message}`,
		)
		backend2Formatter := logging.NewBackendFormatter(backend, format)
		logging.SetBackend(backend2Formatter)
		break
	}

	lvl, err := logging.LogLevel(conf.LogLevel())
	if err != nil {
		panic(err)
	}
	logging.SetLevel(lvl, "")

	LOGGER = logging.NewLog(logging.NewContext())
}
