package conf

import (
	. "github.com/mageddo/dns-proxy-server/log"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/flags"
	"bytes"
	"os"
	"bufio"
	"strings"
	"net"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"io/ioutil"
	"os/exec"
	"syscall"
	"errors"
	"fmt"
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


func GetString(value, defaultValue string) string {

	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func RestoreResolvconfToDefault() error {
	LOGGER.Infof("m=RestoreResolvconfToDefault, status=begin")
	hd := newDNSServerCleanerHandler()
	err := ProcessResolvconf(hd)
	LOGGER.Infof("m=RestoreResolvconfToDefault, status=success, err=%v", err)
	return err
}

func SetMachineDNSServer(serverIP string) error {
	hd := newSetMachineDnsServerHandler(serverIP)
	return ProcessResolvconf(hd)
}

func ProcessResolvconf( handler DnsHandler ) error {

	var newResolvConfBuff bytes.Buffer
	LOGGER.Infof("m=ProcessResolvconf, status=begin")

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
	LOGGER.Infof("m=ProcessResolvconf, status=open-conf-file, file=%s", fileRead.Name())
	scanner := bufio.NewScanner(fileRead)
	for scanner.Scan() {
		line := scanner.Text()
		hasContent = true
		entryType := getDnsEntryType(line)
		if entryType == PROXY {
			foundDnsProxyEntry = true
		}
		LOGGER.Debugf("m=ProcessResolvconf, status=readline, line=%s, type=%s", line,  entryType)
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
	LOGGER.Infof("m=ProcessResolvconf, status=success, buffLength=%d", length)
	return nil
}

func getDNSLine(serverIP string) string {
	return "nameserver " + serverIP + " # dps-entry"
}

func getDnsEntryType(line string) DnsEntry {

	if strings.HasSuffix(line, "# dps-entry") {
		return PROXY
	} else if strings.HasPrefix(line, "# nameserver ") && strings.HasSuffix(line, "# dps-comment") {
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

	ip, err := getCurrentIpAddress()
	LOGGER.Infof("m=SetCurrentDNSServerToMachine, status=begin, ip=%s, err=%v", ip, err)
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

	LOGGER.Infof("m=Lockfile, status=begin, lock=%t, file=%s", lock, file)
	flag := "-i"
	if lock {
		flag = "+i"
	}
	cmd := exec.Command("chattr", flag, file)
	err := cmd.Run()
	if err != nil {
		LOGGER.Warningf("m=Lockfile, status=error-at-execute, lock=%t, file=%s, err=%v", lock, file, err)
		return err
	}
	//bytes, err := cmd.CombinedOutput()

	status := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if status.ExitStatus() != 0 {
		LOGGER.Warningf("m=Lockfile, status=bad-exit-code, lock=%t, file=%s", lock, file)
		return errors.New(fmt.Sprintf("Failed to lock file %d", status.ExitStatus()))
	}
	LOGGER.Infof("m=Lockfile, status=success, lock=%t, file=%s", lock, file)
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

