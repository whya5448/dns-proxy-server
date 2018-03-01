package log

import (
	"github.com/mageddo/go-logging"
	"os"
	"io"
)

var LOGGER logging.Log
var LEVEL string = ""
func init(){
	setup(os.Stdout)
}

func setup(out io.Writer) {

	mode := "dev"
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
	LOGGER = logging.NewLog(logging.NewContext())
	if LEVEL != "" {
		SetLevel(LEVEL)
	}
}

func SetLevel(level string) error {
	LEVEL = level
	lvl, err := logging.LogLevel(level)
	if err != nil {
		return err
	}
	logging.SetLevel(lvl, "")
	return nil
}

func SetOutput(f string) error {
	if f == "console" {
		setup(os.Stdout)
		return nil
	}
	if f == "" {
		SetLevel("CRITICAL")
		return nil
	}

	out, err := os.OpenFile(f, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	setup(out)
	return nil
}
