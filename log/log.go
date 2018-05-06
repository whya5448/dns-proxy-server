package log

import (
	"github.com/mageddo/go-logging"
	"os"
	"io"
)

func init(){
	setup(os.Stdout)
}

func setup(out io.Writer) {
	logging.SetOutput(out)
}

func SetOutput(f string) error {
	if f == "console" {
		setup(os.Stdout)
		return nil
	}

	out, err := os.OpenFile(f, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	setup(out)
	return nil
}
