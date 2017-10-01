package local

import (
	"encoding/json"
	"os"
	log "github.com/mageddo/go-logging"
	"bufio"
	"github.com/mageddo/dns-proxy-server/utils"
	"errors"
	"golang.org/x/net/context"
	"time"
	"fmt"
	"regexp"
	"github.com/mageddo/dns-proxy-server/flags"
	"strings"
)

var confPath string = utils.GetPath(*flags.ConfPath)
var configuration = LocalConfiguration{
	Envs: make([]EnvVo, 0),
	RemoteDnsServers: make([][4]byte, 0),
}

func LoadConfiguration(ctx context.Context){

	logger := log.NewLog(ctx)
	logger.Infof("status=begin, confPath=%s", confPath)

	if _, err := os.Stat(confPath); err == nil {

		logger.Info("status=openingFile")
		f, _ := os.Open(confPath)

		defer f.Close()

		dec := json.NewDecoder(f)
		dec.Decode(&configuration)

		for i := range configuration.Envs {
			env := &configuration.Envs[i]
			for j := range env.Hostnames {
				host := &env.Hostnames[j]
				if host.Id <= 0 {
					logger.Infof("status=without-id, hostname=%s, id=%d", host.Hostname, host.Id)
					host.Id = configuration.nextId()
				}
			}
		}
		logger.Info("status=success")
	}else{
		logger.Info("status=create-new-conf")
		err := os.MkdirAll(confPath[:strings.LastIndex(confPath, "/")], 0755)
		if err != nil {
			logger.Errorf("status=error-to-create-conf-folder, err=%v", err)
			return
		}
		SaveConfiguration(ctx, &configuration)
		logger.Info("status=success")
	}

}
func SaveConfiguration(ctx context.Context, c *LocalConfiguration) {

	t := time.Now()
	logger := log.NewLog(ctx)
	logger.Debugf("status=begin")
	if len(c.Envs) == 0 {
		c.Envs = NewEmptyEnv()
	}

	logger.Debugf("status=save")
	f, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logger.Errorf("status=error-to-create-conf-file, err=%v", err)
		return
	}
	defer f.Close()
	wr := bufio.NewWriter(f)
	defer wr.Flush()
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	err = enc.Encode(c)
	if err != nil {
		logger.Errorf("status=error-to-encode, error=%v", err)
	}
	logger.Infof("status=success, time=%d", utils.DiffMillis(t, time.Now()))

}

func GetConfigurationNoCtx() *LocalConfiguration {
	return GetConfiguration(log.NewContext())
}
func GetConfiguration(ctx context.Context) *LocalConfiguration {
	LoadConfiguration(ctx)
	return &configuration
}


type LocalConfiguration struct {
	/**
	 * The remote servers to ask when, DPS can not solve from docker or local file,
	 * it will try one by one in order, if no one is specified then 8.8.8.8 is used by default
	 * DO NOT call this variable directly, use GetRemoteDnsServers instead
	 */
	RemoteDnsServers [][4]byte `json:"remoteDnsServers"`
	Envs []EnvVo `json:"envs"`
	ActiveEnv string `json:"activeEnv"`
	LastId int `json:"lastId"`

	/// ----
	WebServerPort int `json:"webServerPort"`
	DnsServerPort int `json:"dnsServerPort"`
	DefaultDns *bool `json:"defaultDns"`
}

type EnvVo struct {
	Name string `json:"name"`
	Hostnames []HostnameVo `json:"hostnames"`
}

type HostnameVo struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Ip [4]byte `json:"ip"`
	Ttl int `json:"ttl"`
	Env string `json:"env"` // apenas para o post do rest
}

func (lc *LocalConfiguration) GetEnv(envName string) (*EnvVo, int) {

	for i := range lc.Envs {
		env := &lc.Envs[i]
		if (*env).Name == envName {
			return env, i
		}
	}
	return nil, -1
}

func (foundEnv *EnvVo) AddHostname(ctx context.Context, hostname *HostnameVo) error {

	logger := log.NewLog(ctx)
	logger.Infof("status=begin, env=%s, hostname=%+v", foundEnv.Name, hostname)
	if foundEnv == nil {
		return errors.New("env not found")
	}
	foundHost, _ := foundEnv.GetHostname(hostname.Hostname)
	if foundHost != nil {
		return errors.New(fmt.Sprintf("The host '%s' already exists", hostname.Hostname))
	}

	(*foundEnv).Hostnames = append(foundEnv.Hostnames, *hostname)
	logger.Infof("status=success, foundEnv=%s", foundEnv.Name)
	return nil
}

func (lc *LocalConfiguration) GetActiveEnv() (*EnvVo, int) {
	return lc.GetEnv(lc.ActiveEnv)
}

func(env *EnvVo) GetHostname(hostname string) (*HostnameVo, int) {
	for i := range env.Hostnames {
		host := &env.Hostnames[i]
		if (*host).Hostname == hostname {
			return host, i
		}
	}
	return nil, -1
}

func(env *EnvVo) FindHostnameByName(ctx context.Context, hostname string) *[]HostnameVo {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, hostname=%s", hostname)
	hostList := []HostnameVo{}
	for _, host := range env.Hostnames {
		if matched, _ := regexp.MatchString(fmt.Sprintf(".*%s.*", hostname), host.Hostname); matched {
			hostList = append(hostList, host)
		}
	}
	logger.Infof("status=success, hostname=%s, length=%d", hostname, len(hostList))
	return &hostList
}

func(lc *LocalConfiguration) FindHostnameByNameAndEnv(ctx context.Context, envName, hostname string) (*[]HostnameVo, error) {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, envName=%s, hostname=%s", envName, hostname)
	env,_ := lc.GetEnv(envName)
	if env == nil {
		return nil, errors.New("env not found")
	}
	logger.Infof("status=success, envName=%s, hostname=%s", envName, hostname)
	return env.FindHostnameByName(ctx, hostname), nil
}

func(env *EnvVo) GetHostnameById(id int) (*HostnameVo, int) {
	for i := range env.Hostnames {
		host := &env.Hostnames[i]
		if (*host).Id == id {
			return host, i
		}
	}
	return nil, -1
}

func (lc *LocalConfiguration) AddEnv(ctx context.Context, env EnvVo) error {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, env=%s", env.Name)
	foundEnv, _ := lc.GetEnv(env.Name)
	if foundEnv != nil {
		return errors.New(fmt.Sprintf("The '%s' env already exists", env.Name))
	}
	lc.Envs = append(lc.Envs, env)
	SaveConfiguration(ctx, lc)
	logger.Infof("status=success, env=%s", env.Name)
	return nil
}

func (lc *LocalConfiguration) RemoveEnvByName(ctx context.Context, name string) error {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, env=%s", name)
	env, i := lc.GetEnv(name)
	if env == nil {
		return errors.New(fmt.Sprintf("The env '%s' was not found", name))
	}
	lc.RemoveEnv(ctx, i)
	SaveConfiguration(ctx,lc)
	logger.Infof("status=success, env=%s", name)
	return nil
}

func (lc *LocalConfiguration) RemoveEnv(ctx context.Context, index int){
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, index=%d", index)
	lc.Envs = append(lc.Envs[:index], lc.Envs[index+1:]...)
	SaveConfiguration(ctx,lc)
	logger.Infof("status=success, index=%d", index)
}

func (lc *LocalConfiguration) AddDns(ctx context.Context, dns [4]byte){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers, dns)
	SaveConfiguration(ctx, lc)
}

func (lc *LocalConfiguration) RemoveDns(ctx context.Context, index int){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers[:index], lc.RemoteDnsServers[index+1:]...)
	SaveConfiguration(ctx, lc)
}


func (lc *LocalConfiguration) AddHostname(ctx context.Context, envName string, hostname HostnameVo) error {
	logger := log.NewLog(ctx)
	hostname.Id = lc.nextId()
	logger.Infof("status=begin, evnName=%s, hostname=%+v", envName, hostname)
	foundEnv, _ := lc.GetEnv(envName)
	if foundEnv == nil {
		return errors.New("env not found")
	}
	err := foundEnv.AddHostname(ctx, &hostname)
	if err != nil {
		return err
	}
	SaveConfiguration(ctx, lc)
	logger.Infof("status=success")
	return nil
}

func (lc *LocalConfiguration) nextId() int {
	lc.LastId++;
	return lc.LastId;
}

func (lc *LocalConfiguration) UpdateHostname(ctx context.Context, envName string, hostname HostnameVo) error {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, evnName=%s, hostname=%+v", envName, hostname)
	env, _ := lc.GetEnv(envName)
	if(env == nil){
		return errors.New("env not found")
	}

	err := env.UpdateHostname(hostname)
	if err != nil {
		return err
	}

	SaveConfiguration(ctx, lc)
	logger.Infof("status=success")
	return nil
}

func (env *EnvVo) UpdateHostname(hostname HostnameVo) error {

	foundHostname, _ := env.GetHostnameById(hostname.Id)
	if foundHostname == nil {
		return errors.New("not hostname found")
	}
	foundHostname.Ip = hostname.Ip;
	foundHostname.Ttl = hostname.Ttl;
	foundHostname.Hostname = hostname.Hostname;

	return nil
}

func (lc *LocalConfiguration) RemoveHostnameByEnvAndHostname(ctx context.Context, envName string, hostname string) error {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, envName=%s, hostname=%s", envName, hostname)
	env, envIndex := lc.GetEnv(envName)
	if envIndex == -1 {
		return errors.New("env not found")
	}
	host, hostIndex := env.GetHostname(hostname)
	if host == nil {
		return errors.New("hostname not found")
	}
	lc.RemoveHostname(ctx, envIndex, hostIndex)
	logger.Infof("status=success, envName=%s, hostname=%s", envName, hostname)
	return nil
}

func (lc *LocalConfiguration) RemoveHostname(ctx context.Context, envIndex int, hostIndex int){

	logger := log.NewLog(ctx)
	logger.Infof("status=begin, envIndex=%d, hostIndex=%d", envIndex, hostIndex)
	env := &lc.Envs[envIndex];
	(*env).Hostnames = append((*env).Hostnames[:hostIndex], (*env).Hostnames[hostIndex+1:]...)
	SaveConfiguration(ctx, lc)
	logger.Infof("status=success, envIndex=%d, hostIndex=%d", envIndex, hostIndex)

}

func (lc *LocalConfiguration) SetActiveEnv(ctx context.Context, env EnvVo) error {
	logger := log.NewLog(ctx)
	logger.Infof("status=begin, envActive=%s", env.Name)
	foundEnv, _ := lc.GetEnv(env.Name)
	if foundEnv == nil {
		logger.Warningf("status=env-not-found, envName=%s", env.Name)
		return errors.New("Env not found: " + env.Name)
	}
	lc.ActiveEnv = env.Name
	SaveConfiguration(ctx, lc)
	logger.Infof("status=success")
	return nil
}

func NewEmptyEnv() []EnvVo {
	return []EnvVo{{Hostnames:[]HostnameVo{}, Name:""}}
}

func (lc *LocalConfiguration) GetRemoteServers(ctx context.Context) [][4]byte {
	if len(lc.RemoteDnsServers) == 0 {
		lc.RemoteDnsServers = append(lc.RemoteDnsServers, [4]byte{8, 8, 8, 8})
		logger := log.NewLog(ctx)
		logger.Infof("status=put-default-server")
	}
	return lc.RemoteDnsServers
}
