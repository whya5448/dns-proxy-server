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
	"github.com/mageddo/dns-proxy-server/conf"
	"bufio"
	"bytes"
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
			resp.Compress = conf.Compress()
			respWriter.WriteMsg(resp)
			break
		}

		logger.Warningf("status=not-resolved, solver=%s, err=%v", solverID, err)

	}

}

func serve(net, name, secret string, logger *log.IdLogger) {
	port := fmt.Sprintf(":%d", conf.DnsServerPort())
	logger.Infof("status=begin, port=%d", conf.DnsServerPort())
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
	if conf.Tsig() != "" {
		a := strings.SplitN(conf.Tsig(), ":", 2)
		name, secret = dns.Fqdn(a[0]), a[1] // fqdn the name, which everybody forgets...
	}
	if conf.CpuProfile() != "" {
		f, err := os.Create(conf.CpuProfile())
		if err != nil {
			logger.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dns.HandleFunc(".", handleQuestion)

	local.GetConfiguration(context)

	go docker.HandleDockerEvents()
	go serve("tcp", name, secret, logger)
	go serve("udp", name, secret, logger)
	go func(){
		webPort := conf.WebServerPort();
		logger.Infof("status=web-server-starting, port=%d", webPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", webPort), nil); err != nil {
			logger.Errorf("status=failed-start-web-server, err=%v, port=%d", err, webPort)
			os.Exit(3)
		}else{
			logger.Infof("status=web-server-started, port=%d", webPort)
		}
	}()

	controller.MapRequests()

	var buffer bytes.Buffer

	fmt.Println(buffer.String())
	file, err := os.Open("/path/to/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		if strings.HasSuffix(line, "# dns-proxy-server") { // this line is dns proxy server nameserver entry
			ip := "" // TODO recuperar o IP da maquina aqui
			buffer.WriteString(ip + "# dns-proxy-server")
		}else if strings.HasPrefix(line, "#") { // linha comentada
			buffer.WriteString(line)
		} else if strings.HasPrefix(line, "nameserver") {
			buffer.WriteString("# " + line)
		} else {
			//
		}
		buffer.WriteString("\n")


		buffer.WriteString()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)
}
