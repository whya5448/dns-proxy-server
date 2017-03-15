package local

import (
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
	"github.com/mageddo/dns-proxy-server/utils"
	"errors"
	"golang.org/x/net/context"
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
	logger.Infof("m=SaveConfiguration, status=begin, configuration=%+v", c)
	if len(c.Envs) == 0 {
		c.Envs = NewEmptyEnv()
	}

	js,_ := json.Marshal(&c)
	logger.Infof("m=SaveConfiguration, status=save, data=%s", js)

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
	err = enc.Encode(c)
	if err != nil {
		logger.Errorf("status=error-to-encode, error=%v", err)
	}

	logger.Infof("m=SaveConfiguration, status=success")

}

func GetConfiguration(ctx context.Context) LocalConfiguration {
	LoadConfiguration(ctx)
	return configuration
}


type LocalConfiguration struct {
	RemoteDnsServers [][4]byte `json:"remoteDnsServers"`
	Envs []EnvVo `json:"envs"`
	ActiveEnv string `json:"activeEnv"`
}

type EnvVo struct {
	Name string `json:"name"`
	Hostnames []HostnameVo `json:"hostnames"`
}

type HostnameVo struct {
	Hostname string `json:"hostname"`
	Ip [4]byte `json:"ip"`
	Ttl int `json:"ttl"`
	Env string `json:"env"` // apenas para o post do rest
}

func (lc *LocalConfiguration) GetEnv(envName string) (*EnvVo) {

	for i := range lc.Envs {
		env := &lc.Envs[i]
		if (*env).Name == envName {
			return env
		}
	}
	return nil
}

func (lc *LocalConfiguration) AddHostnameToEnv(env string, hostname *HostnameVo) error {
	log.Logger.Infof("m=AddHostnameToEnv, status=begin, env=%+v, hostname=%+v", env, hostname)
	foundEnv := lc.GetEnv(env)
	if foundEnv == nil {
		return errors.New("env not found")
	}
	(*foundEnv).Hostnames = append(foundEnv.Hostnames, *hostname)
	log.Logger.Infof("m=AddHostnameToEnv, status=success, lc=%+v, foundEnv=%+v, hostnames=%+v", lc, foundEnv, lc.Envs[0].Hostnames)
	return nil
}

func (lc *LocalConfiguration) GetActiveEnv() *EnvVo {
	return lc.GetEnv(lc.ActiveEnv)
}

func(env *EnvVo) GetHostname(hostname string) *HostnameVo {
	for _, host := range env.Hostnames {
		if host.Hostname == hostname {
			return &host
		}
	}
	return nil
}

func (lc *LocalConfiguration) AddEnv(ctx context.Context, env EnvVo){
	lc.Envs = append(lc.Envs, env)
	SaveConfiguration(ctx, lc)
}

func (lc *LocalConfiguration) RemoveEnv(ctx context.Context,index int){
	lc.Envs = append(lc.Envs[:index], lc.Envs[index+1:]...)
	SaveConfiguration(ctx,lc)
}

func (lc *LocalConfiguration) AddDns(ctx context.Context, dns [4]byte){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers, dns)
	SaveConfiguration(ctx, lc)
}

func (lc *LocalConfiguration) RemoveDns(ctx context.Context, index int){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers[:index], lc.RemoteDnsServers[index+1:]...)
	SaveConfiguration(ctx, lc)
}


func AddHostname(ctx context.Context, envName string, hostname HostnameVo) error {
	log.Logger.Infof("m=AddHostname, status=begin, evnName=%s, hostname=%+v", envName, hostname)
	err := configuration.AddHostnameToEnv(envName, &hostname)
	if err != nil {
		return err
	}
	SaveConfiguration(ctx, &configuration)
	log.Logger.Infof("m=AddHostname, status=success, configuration=%+v", configuration)
	return nil
}

func RemoveHostname(ctx context.Context, envIndex int, hostIndex int){
	env := configuration.Envs[envIndex];
	t := append(env.Hostnames[:hostIndex], env.Hostnames[hostIndex+1:]...)
	env.Hostnames = t
	SaveConfiguration(ctx, &configuration)
}

func NewEmptyEnv() []EnvVo {
	return []EnvVo{{Hostnames:[]HostnameVo{}, Name:""}}
}