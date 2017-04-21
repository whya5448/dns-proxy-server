package log

import (
	"os"
	"github.com/mageddo/go-logging"
	"strings"
	"io"
	"fmt"
	"golang.org/x/net/context"
	"bytes"
)
var Logger *logging.Logger
var out io.Writer = os.Stdout
var mode string = strings.ToUpper(os.Getenv("MG_ENV"))

func init(){
	setup()
}

func setup() {

	Logger = logging.MustGetLogger("main")
	mode := getMode()
	backend := logging.NewLogBackend(getOut(), "", 0)

	// setando o log dependendo do ambiente
	switch mode {
	case "DEV":
		format := logging.MustStringFormatter(
			`%{color}%{time:06-01-02 15:04:05.000} %{level:.3s} %{color:reset}%{message}`,
		)
		backend2Formatter := logging.NewBackendFormatter(backend, format)
		logging.SetBackend(backend2Formatter)
		break;
	case "PROD":
		format := logging.MustStringFormatter(
			`%{time:06-01-02 15:04:05} %{level:.3s} %{message}`,
		)
		leveledBackend := logging.AddModuleLevel(logging.NewBackendFormatter(backend, format));
		leveledBackend.SetLevel(logging.INFO, "")
		logging.SetBackend(leveledBackend)
		break;
	}

}

func getOut() io.Writer {
	return out
}
func getMode() string {
	if(len(mode) == 0){
		return "DEV"
	}
	return mode
}

type IdLogger int

func GetLogger(ctx context.Context) *IdLogger {
	x := IdLogger(ctx.Value(LoggerID).(int))
	return &x
}

func (l *IdLogger) Critical(args ...interface{}) {
	Logger.Critical(getArgs(l, args...)...)
}
func (l *IdLogger) Criticalf(format string, args ...interface{}) {
	Logger.Criticalf(getIdConcat(l, format, -1), args...)
}
func (l *IdLogger) Debug(args ...interface{}) {
	Logger.Debug(getArgs(l, args...)...)
}
func (l *IdLogger) Debugf(format string, args ...interface{}) {
	Logger.Debugf(getIdConcat(l, format, -1), args...)
}
func (l *IdLogger) Error(args ...interface{}) {
	Logger.Error(getArgs(l, args...)...)
}
func (l *IdLogger) Errorf(format string, args ...interface{}) {
	Logger.Errorf(getIdConcat(l, format, -1), args...)
}
func (l *IdLogger) Fatal(args ...interface{}) {
	Logger.Fatal(getArgs(l, args...)...)
	os.Exit(1)
}
func (l *IdLogger) Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(getIdConcat(l, format, -1), args...)
	os.Exit(1)
}
func (l *IdLogger) Info(args ...interface{}) {
	Logger.Info(getArgs(l, args...)...)
}
func (l *IdLogger) Infof(format string, args ...interface{}) {
	Logger.Infof(getIdConcat(l, format, -1), args...)

}

func (l *IdLogger) Notice(args ...interface{}) {
	Logger.Notice(getArgs(l, args...)...)
}
func (l *IdLogger) Noticef(format string, args ...interface{}) {
	Logger.Noticef(getIdConcat(l, format, -1), args...)
}
func (l *IdLogger) Panic(args ...interface{}) {
	Logger.Panic(args)
	panic(fmt.Sprint(args...))
}
func (l *IdLogger) Panicf(format string, args ...interface{}) {
	Logger.Panicf(getIdConcat(l, format, -1), args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *IdLogger) Warning(args ...interface{}) {
	Logger.Warning(getArgs(l, args...)...)
}
func (l *IdLogger) Warningf(format string, args ...interface{}) {
	Logger.Warningf(getIdConcat(l, format, -1), args...)
}

func getFormat(format string) string{
	var buffer bytes.Buffer
	buffer.WriteString("ctx=%d, m=")
	buffer.WriteString(GetCallerFunctionNameSkippingAnnonymous(2))
	buffer.WriteString(" ")
	buffer.WriteString(format)
	return buffer.String()
}


func getArgs(l *IdLogger, args ...interface{}) []interface{} {
	var ar  []interface{}
	ar = append(ar, strings.Trim(getId(l), " "))
	ar = append(ar, args...)
	return ar;
}

func getId(l *IdLogger) string {
	return getIdConcat(l, "", 4)
}

func getIdConcat(l *IdLogger, append string, level int ) string {
	if level == -1{
		level = 2
	}
	var buffer bytes.Buffer
	buffer.WriteString("ctx=")
	buffer.WriteString(fmt.Sprintf("%d, ", *l))
	buffer.WriteString("m=")
	buffer.WriteString(GetCallerFunctionNameSkippingAnnonymous(level))
	buffer.WriteString(" ")
	buffer.WriteString(append)
	return buffer.String()
}