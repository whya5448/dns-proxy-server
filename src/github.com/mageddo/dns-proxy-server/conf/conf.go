package conf

import (
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"bytes"
	"os"
	"bufio"
	"strings"
	"net"
	"github.com/mageddo/log"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"io/ioutil"
	"os/exec"
	"syscall"
	"errors"
	"fmt"
	"github.com/mageddo/dns-proxy-server/utils"
)

type DnsEntry string

const(
	COMMENT DnsEntry = "COMMENT"
	COMMENTED_SERVER DnsEntry = "COMMENTED_SERVER"
	SERVER DnsEntry = "SERVER"
	PROXY DnsEntry = "PROXY"
	ELSE DnsEntry = "ELSE"
)

type DnsHandler interface {
	process(line string, entryType DnsEntry) *string
	afterProcess(hasContent bool, foundDnsProxy bool) *string
}

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
	port := local.GetConfigurationNoCtx().WebServerPort
	if port <= 0 {
		return *flags.WebServerPort
	}
	return port
}

func DnsServerPort() int {
	port := local.GetConfigurationNoCtx().DnsServerPort
	if port <= 0 {
	return *flags.DnsServerPort
	}
	return port
}

func SetupResolvConf() bool {

	defaultDns := local.GetConfigurationNoCtx().DefaultDns
	if defaultDns == nil {
		return *flags.SetupResolvconf
	}
	return *defaultDns

}

func ConfPath() string {
	return *flags.ConfPath
}

func SetupService() bool {
	return *flags.SetupService
}

func GetString(value, defaultValue string) string {

	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func RestoreResolvconfToDefault() error {
	hd := newDNSServerCleanerHandler()
	return ProcessResolvconf(hd)
}

func SetMachineDNSServer(serverIP string) error {
	hd := newSetMachineDnsServerHandler(serverIP)
	return ProcessResolvconf(hd)
}

func ProcessResolvconf( handler DnsHandler ) error {

	var newResolvConfBuff bytes.Buffer
	log.Logger.Infof("m=ProcessResolvconf, status=begin")

	resolvconf := getResolvConf()
	fileRead, err := os.Open(resolvconf)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	var (
		hasContent = false
		foundDnsProxyEntry = false
	)
	log.Logger.Infof("m=ProcessResolvconf, status=open-conf-file, file=%s", fileRead.Name())
	scanner := bufio.NewScanner(fileRead)
	for scanner.Scan() {
		line := scanner.Text()
		hasContent = true
		entryType := getDnsEntryType(line)
		log.Logger.Debugf("m=ProcessResolvconf, status=readline, line=%s, type=%s", line,  entryType)
		if r := handler.process(line, entryType); r != nil {
			newResolvConfBuff.WriteString(*r)
			newResolvConfBuff.WriteByte('\n')
		}
	}
	if r := handler.afterProcess(hasContent, foundDnsProxyEntry); r != nil {
		newResolvConfBuff.WriteString(*r)
		newResolvConfBuff.WriteByte('\n')
	}

	stats, _ := fileRead.Stat()
	length := newResolvConfBuff.Len()
	err = ioutil.WriteFile(resolvconf, newResolvConfBuff.Bytes(), stats.Mode())
	if err != nil {
		return err
	}
	log.Logger.Infof("m=ProcessResolvconf, status=success, buffLength=%d", length)
	return nil
}

func getDNSLine(serverIP string) string {
	return "nameserver " + serverIP + " # dns-proxy-server"
}

func getDnsEntryType(line string) DnsEntry {

	if strings.HasSuffix(line, "# dns-proxy-server") {
		return PROXY
	} else if strings.HasPrefix(line, "# nameserver ") {
		return COMMENTED_SERVER
	} else if strings.HasPrefix(line, "#") {
		return COMMENT
	} else if strings.HasPrefix(line, "nameserver") {
		return SERVER
	} else {
		return ELSE
	}
}

func SetCurrentDNSServerToMachineAndLockIt() error {

	err := SetCurrentDNSServerToMachine()
	if err != nil {
		return err
	}
	return LockResolvConf()

}

func SetCurrentDNSServerToMachine() error {

	log.Logger.Infof("m=SetCurrentDNSServerToMachine, status=begin")
	ip, err := getCurrentIpAddress()
	if err != nil {
		return err
	}
	return SetMachineDNSServer(ip)
}

func LockResolvConf() error {
	return LockFile(true, getResolvConf())
}

func UnlockResolvConf() error {
	return LockFile(true, getResolvConf())
}

func LockFile(lock bool, file string) error {

	log.Logger.Infof("m=Lockfile, status=begin, lock=%t, file=%s", lock, file)
	flag := "-i"
	if lock {
		flag = "+i"
	}
	cmd := exec.Command("chattr", flag, file)
	err := cmd.Run()
	if err != nil {
		log.Logger.Warningf("m=Lockfile, status=error-at-execute, lock=%t, file=%s, err=%v", lock, file, err)
		return err
	}
	//bytes, err := cmd.CombinedOutput()

	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if status.ExitStatus() != 0 {
		log.Logger.Warningf("m=Lockfile, status=bad-exit-code, lock=%t, file=%s", lock, file)
		return errors.New(fmt.Sprintf("Failed to lock file %d", status.ExitStatus()))
	}
	log.Logger.Infof("m=Lockfile, status=success, lock=%t, file=%s", lock, file)
	return nil

}

func getCurrentIpAddress() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ip := addr.String()
		if strings.Contains(ip, "/") {
			if !strings.HasPrefix(ip, "127") {
				return ip[:strings.Index(ip, "/")], nil
			}
		}
	}
	return "", nil

}

func getResolvConf() string {
	return GetString(os.Getenv(env.MG_RESOLVCONF), "/etc/resolv.conf")
}

func ConfigSetupService(){

	log.Logger.Infof("m=ConfigSetupService, status=begin")
	err := utils.Copy("dns-proxy-server", "/etc/init.d/dns-proxy-server")
	//err = ioutil.WriteFile(resolvconf, newResolvConfBuff.Bytes(), stats.Mode())
	if err != nil {
		log.Logger.Fatalf("status=error-copy-service, msg=%s", err.Error())
	}
	err = utils.Copy("docker-compose.yml", "/etc/init.d/dns-proxy-server.yml")
	if err != nil {
		log.Logger.Fatalf("status=error-copy-yml, msg=%s", err.Error())
	}
	err = utils.Exec("/usr/sbin/update-rc.d", "dns-proxy-server", "defaults")
	if err != nil {
		log.Logger.Fatalf("status=fatal-install-service, msg=%s", err.Error())
	}
	err = utils.Exec("/usr/sbin/service", "dns-proxy-server", "start")
	if err != nil {
		log.Logger.Fatalf("status=start-service, msg=%s", err.Error())
	}
	log.Logger.Infof("m=ConfigSetupService, status=success")

}