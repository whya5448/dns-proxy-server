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
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/dns-proxy-server/proxy"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	compress   = flag.Bool("compress", false, "compress replies")
	tsig       = flag.String("tsig", "", "use MD5 hmac tsig: keyname:base64")
)

func handleQuestion(respWriter dns.ResponseWriter, reqMsg *dns.Msg) {

	defer func() {
		err := recover()
		if err != nil {
			log.Logger.Errorf("M=handleReflect, status=error, error=%v", err)
		}
	}()

	var questionName string
	questionsQtd := len(reqMsg.Question)
	if questionsQtd != 0 {
		questionName = reqMsg.Question[0].Name
	}	else {
		questionName = "null"
	}

	log.Logger.Infof("m=handleReflect, questions=%d, 1stQuestion=%s", questionsQtd, questionName)


	// loading the solvers and try to solve the hostname in that order
	solvers := []proxy.DnsSolver{proxy.LocalDnsSolver{}, proxy.DockerDnsSolver{}, proxy.RemoteDnsSolver{}}
	for _, solver := range solvers {

		// loop through questions
		answer := solver.Solve(questionName)

		answer.SetReply(reqMsg)
		answer.Compress = *compress
		respWriter.WriteMsg(answer);

	}


	resp := utils.SolveName(questionName)
	resp.SetReply(reqMsg)
	resp.Compress = *compress


	var firstAnswer dns.RR
	if len(resp.Answer) != 0 {
		firstAnswer = resp.Answer[0]
	}

	log.Logger.Infof("m=handleReflect, resp=%v", firstAnswer)
	respWriter.WriteMsg(resp)

}

func serve(net, name, secret string) {
	switch name {
	case "":
		server := &dns.Server{Addr: ":53", Net: net, TsigSecret: nil}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
		}
	default:
		server := &dns.Server{Addr: ":53", Net: net, TsigSecret: map[string]string{name: secret}}
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
	go serve("tcp", name, secret)
	go serve("udp", name, secret)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)
}
