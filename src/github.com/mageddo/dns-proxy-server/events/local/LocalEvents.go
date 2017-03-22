package local

import (
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
	"github.com/mageddo/dns-proxy-server/utils"
	"errors"
	"golang.org/x/net/context"
	"time"
	"fmt"
)

var confPath string = utils.GetPath("conf/config.json")
var configuration = LocalConfiguration{
	Envs: make([]EnvVo, 0),
	RemoteDnsServers: make([][4]byte, 0),
}

func init(){
	if len(os.Args) > 2 {
		confPath = utils.GetPath(os.Args[2]);
		log.Logger.Infof("m=init, status=changed-confpath, confpath=%s", utils.GetPath(confPath))
	}
	LoadConfiguration(log.GetContext())

}

func LoadConfiguration(ctx context.Context){

	logger := log.GetLogger(ctx)
	logger.Infof("status=begin, confPath=%s", confPath)

	if _, err := os.Stat(confPath); err == nil {
		logger.Info("status=openingFile")
		f, _ := os.Open(confPath)

		defer func(){
			f.Close()
		}()

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

		SaveConfiguration(ctx, &configuration)
		logger.Info("status=success")
	}else{
		logger.Info("status=create-new-conf")
		err := os.MkdirAll(confPath, 0755)
		if err != nil {
			logger.Errorf("status=error-to-create-conf-folder, err=%v", err)
			return
		}
		SaveConfiguration(ctx, &configuration)
		logger.Info("status=success")
	}

}
func SaveConfiguration(ctx context.Context, c *LocalConfiguration) {
	logger := log.GetLogger(ctx)

	logger.Infof("m=SaveConfiguration, status=begin, time=%s", time.Now())
	if len(c.Envs) == 0 {
		c.Envs = NewEmptyEnv()
	}

	logger.Infof("m=SaveConfiguration, status=save")
	f, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	defer func(){
		f.Close()
	}()
	if err != nil {
		logger.Errorf("status=error-to-create-conf-file, err=%v", err)
		return
	}
	wr := bufio.NewWriter(f)
	defer func(){
		wr.Flush()
	}()
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	err = enc.Encode(c)
	if err != nil {
		logger.Errorf("status=error-to-encode, error=%v", err)
	}

	logger.Infof("m=SaveConfiguration, status=success")

}

func GetConfiguration(ctx context.Context) *LocalConfiguration {
	LoadConfiguration(ctx)
	return &configuration
}


type LocalConfiguration struct {
	RemoteDnsServers [][4]byte `json:"remoteDnsServers"`
	Envs []EnvVo `json:"envs"`
	ActiveEnv string `json:"activeEnv"`
	LastId int `json:"lastId"`
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

func (lc *LocalConfiguration) AddHostnameToEnv(ctx context.Context, env string, hostname *HostnameVo) error {

	logger := log.GetLogger(ctx)
	logger.Infof("m=AddHostnameToEnv, status=begin, env=%+v, hostname=%+v", env, hostname)
	foundEnv, _ := lc.GetEnv(env)
	if foundEnv == nil {
		return errors.New("env not found")
	}
	foundHost, _ := foundEnv.GetHostname(hostname.Hostname)
	if foundHost != nil {
		return errors.New(fmt.Sprintf("The host '%s' already exists", hostname.Hostname))
	}

	(*foundEnv).Hostnames = append(foundEnv.Hostnames, *hostname)
	logger.Infof("status=success, foundEnv=%s, hostnames=%d", foundEnv.Name, len(lc.Envs))
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
	logger := log.GetLogger(ctx)
	logger.Infof("m=AddEnv, status=begin, env=%s", env.Name)
	foundEnv, _ := lc.GetEnv(env.Name)
	if foundEnv != nil {
		return errors.New(fmt.Sprintf("The '%s' env already exists", env.Name))
	}
	lc.Envs = append(lc.Envs, env)
	SaveConfiguration(ctx, lc)
	logger.Infof("m=AddEnv, status=success, env=%s", env.Name)
	return nil
}

func (lc *LocalConfiguration) RemoveEnvByName(ctx context.Context, name string) error {
	logger := log.GetLogger(ctx)
	logger.Infof("m=RemoveEnvByName, status=begin, env=%s", name)
	env, i := lc.GetEnv(name)
	if env == nil {
		return errors.New(fmt.Sprintf("The env '%s' was not found", name))
	}
	lc.RemoveEnv(ctx, i)
	SaveConfiguration(ctx,lc)
	logger.Infof("m=RemoveEnvByName, status=success, env=%s", name)
	return nil
}

func (lc *LocalConfiguration) RemoveEnv(ctx context.Context, index int){
	logger := log.GetLogger(ctx)
	logger.Infof("m=RemoveEnv, status=begin, index=%d", index)
	lc.Envs = append(lc.Envs[:index], lc.Envs[index+1:]...)
	SaveConfiguration(ctx,lc)
	logger.Infof("m=RemoveEnv, status=success, index=%d", index)
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
	logger := log.GetLogger(ctx)
	hostname.Id = lc.nextId()
	logger.Infof("status=begin, evnName=%s, hostname=%+v", envName, hostname)
	err := lc.AddHostnameToEnv(ctx, envName, &hostname)
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
	logger := log.GetLogger(ctx)
	logger.Infof("status=begin, evnName=%s, hostname=%+v", envName, hostname)
	env, _ := lc.GetEnv(envName)
	if(env == nil){
		return errors.New("env not found")
	}
	foundHostName, _ := env.GetHostnameById(hostname.Id)
	if foundHostName == nil {
		return errors.New("not hostname found")
	}

	foundHostName.Ip = hostname.Ip;
	foundHostName.Ttl = hostname.Ttl;
	foundHostName.Hostname = hostname.Hostname;

	SaveConfiguration(ctx, lc)
	logger.Infof("status=success")
	return nil
}

func (lc *LocalConfiguration) RemoveHostnameByEnvAndHostname(ctx context.Context, envName string, hostname string) error {
	logger := log.GetLogger(ctx)
	logger.Infof("m=RemoveHostnameByEnvAndHostname, status=begin, envName=%s, hostname=%s", envName, hostname)
	env, envIndex := lc.GetEnv(envName)
	if envIndex == -1 {
		return errors.New("env not found")
	}
	host, hostIndex := env.GetHostname(hostname)
	if host == nil {
		return errors.New("hostname not found")
	}
	lc.RemoveHostname(ctx, envIndex, hostIndex)
	logger.Infof("m=RemoveHostnameByEnvAndHostname, status=success, envName=%s, hostname=%s", envName, hostname)
	return nil
}

func (lc *LocalConfiguration) RemoveHostname(ctx context.Context, envIndex int, hostIndex int){
	logger := log.GetLogger(ctx)
	logger.Infof("m=RemoveHostname, status=begin, envIndex=%d, hostIndex=%d", envIndex, hostIndex)
	env := &lc.Envs[envIndex];
	(*env).Hostnames = append((*env).Hostnames[:hostIndex], (*env).Hostnames[hostIndex+1:]...)
	SaveConfiguration(ctx, lc)
	logger.Infof("m=RemoveHostname, status=success, envIndex=%d, hostIndex=%d", envIndex, hostIndex)
}

func (lc *LocalConfiguration) SetActiveEnv(ctx context.Context, env EnvVo) error {
	logger := log.GetLogger(ctx)
	logger.Infof("m=SetActiveEnv, status=begin, envActive=%s", env.Name)
	foundEnv, _ := lc.GetEnv(env.Name)
	if foundEnv == nil {
		logger.Warningf("m=SetActiveEnv, status=env-not-found, envName=%s", env.Name)
		return errors.New("Env not found: " + env.Name)
	}
	lc.ActiveEnv = env.Name
	SaveConfiguration(ctx, lc)
	logger.Infof("m=SetActiveEnv, status=success")
	return nil
}

func NewEmptyEnv() []EnvVo {
	return []EnvVo{{Hostnames:[]HostnameVo{}, Name:""}}
}