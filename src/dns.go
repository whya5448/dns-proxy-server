package main

import (
	"flag"
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
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	compress   = flag.Bool("compress", false, "compress replies")
	tsig       = flag.String("tsig", "", "use MD5 hmac tsig: keyname:base64")
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
	solvers := []proxy.DnsSolver{/*proxy.LocalDnsSolver{},*/ proxy.DockerDnsSolver{}, proxy.RemoteDnsSolver{}}
	for _, solver := range solvers {

		solverID := reflect.TypeOf(solver).Name()
		// loop through questions
		resp, err := solver.Solve(firstQuestion)
		if err == nil {

			var firstAnswer dns.RR
			answerLenth := len(resp.Answer)
			if answerLenth != 0 {
				firstAnswer = resp.Answer[0]
			}
			logger.Infof("status=resolved, solver=%s, alength=%d, answer=%v", solverID, answerLenth, firstAnswer)

			resp.SetReply(reqMsg)
			resp.Compress = *compress
			respWriter.WriteMsg(resp)
			break
		}

		logger.Warningf("status=not-resolved, solver=%s, err=%v", solverID, err)

	}

}

const serverPort = 53

func serve(net, name, secret string) {
	var port string = fmt.Sprintf(":%d", serverPort)
	switch name {
	case "":
		server := &dns.Server{Addr: port, Net: net, TsigSecret: nil}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
		}
	default:
		server := &dns.Server{Addr: port, Net: net, TsigSecret: map[string]string{name: secret}}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
		}
	}
}

func main() {
	var name, secret string
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()
	if *tsig != "" {
		a := strings.SplitN(*tsig, ":", 2)
		name, secret = dns.Fqdn(a[0]), a[1] // fqdn the name, which everybody forgets...
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Logger.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dns.HandleFunc(".", handleQuestion)

	local.GetConfiguration()

	go docker.HandleDockerEvents()
	go serve("tcp", name, secret)
	go serve("udp", name, secret)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)
}
