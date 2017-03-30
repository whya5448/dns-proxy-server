package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
	"github.com/miekg/dns"
	"github.com/mageddo/log"
	"github.com/mageddo/dns-proxy-server/proxy"
	"reflect"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/events/docker"
	"net/http"
	"github.com/mageddo/dns-proxy-server/controller"
	"github.com/mageddo/dns-proxy-server/flags"
)

func handleQuestion(respWriter dns.ResponseWriter, reqMsg *dns.Msg) {

	ctx := log.GetContext()
	logger := log.GetLogger(ctx)

	defer func() {
		err := recover()
		if err != nil {
			logger.Errorf("status=error, error=%v", err)
		}
	}()

	var firstQuestion dns.Question
	questionsQtd := len(reqMsg.Question)
	if questionsQtd != 0 {
		firstQuestion = reqMsg.Question[0]
	}else{
		logger.Error("status=question-is-nil")
		return
	}

	logger.Infof("status=begin, reqId=%d, questions=%d, question=%s, type=%s", reqMsg.Id,
		questionsQtd, firstQuestion.Name, utils.DnsQTypeCodeToName(firstQuestion.Qtype))

	// loading the solvers and try to solve the hostname in that order
	solvers := []proxy.DnsSolver{proxy.LocalDnsSolver{}, proxy.DockerDnsSolver{}, proxy.RemoteDnsSolver{}}
	for _, solver := range solvers {

		solverID := reflect.TypeOf(solver).Name()
		logger.Infof("status=begin, solver=%s", solverID)
		// loop through questions
		resp, err := solver.Solve(ctx, firstQuestion)
		if err == nil {

			var firstAnswer dns.RR
			answerLenth := len(resp.Answer)

			logger.Infof("status=answer-found, solver=%s, length=%d", solverID, answerLenth)
			if answerLenth != 0 {
				firstAnswer = resp.Answer[0]
			}
			logger.Infof("status=resolved, solver=%s, alength=%d, answer=%v", solverID, answerLenth, firstAnswer)

			resp.SetReply(reqMsg)
			resp.Compress = *flags.Compress
			respWriter.WriteMsg(resp)
			break
		}

		logger.Warningf("status=not-resolved, solver=%s, err=%v", solverID, err)

	}

}

func serve(net, name, secret string) {
	port := fmt.Sprintf(":%d", *flags.DnsServerPort)
	switch name {
	case "":
		server := &dns.Server{Addr: port, Net: net, TsigSecret: nil}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the %s server: %s\n", net, err.Error())
			os.Exit(2)
		}
	default:
		server := &dns.Server{Addr: port, Net: net, TsigSecret: map[string]string{name: secret}}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the %s server: %s\n", net, err.Error())
			os.Exit(2)
		}
	}
}

func main() {

	context := log.GetContext()
	logger := log.GetLogger(context)

	var name, secret string
	if *flags.Tsig != "" {
		a := strings.SplitN(*flags.Tsig, ":", 2)
		name, secret = dns.Fqdn(a[0]), a[1] // fqdn the name, which everybody forgets...
	}
	if *flags.Cpuprofile != "" {
		f, err := os.Create(*flags.Cpuprofile)
		if err != nil {
			logger.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dns.HandleFunc(".", handleQuestion)

	local.GetConfiguration(context)

	go docker.HandleDockerEvents()
	go serve("tcp", name, secret)
	go serve("udp", name, secret)
	go func(){
		webPort := *flags.WebServerPort;
		logger.Infof("status=web-server-starting, port=%d", webPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", webPort), nil); err != nil {
			logger.Errorf("status=failed-start-web-server, err=%v, port=%d", err, webPort)
			os.Exit(3)
		}else{
			logger.Infof("status=web-server-started, port=%d", webPort)
		}
	}()

	controller.MapRequests()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)
}
