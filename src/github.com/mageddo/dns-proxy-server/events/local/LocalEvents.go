package local

import (
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
	"github.com/mageddo/dns-proxy-server/utils"
)

var confPath string = "conf/config.json"
var configuration = LocalConfiguration{
	Envs: make([]EnvVo, 0),
	RemoteDnsServers: make([][4]byte, 0),
}

func init(){
	if len(os.Args) > 2 {
		confPath = os.Args[2];
		log.Logger.Infof("m=init, status=changed-confpath, confpath=%s", utils.GetPath(confPath))
	}

}

func LoadConfiguration(){

	if _, err := os.Stat(confPath); err == nil {

		f, _ := os.Open(confPath)

		defer func(){
			f.Close()
		}()

		dec := json.NewDecoder(f)
		dec.Decode(&configuration)
		SaveConfiguration(&configuration)

	}else{
		err := os.MkdirAll("conf", 0755)
		if err != nil {
			log.Logger.Errorf("status=error-to-create-conf-folder, err=%v", err)
			return
		}
		SaveConfiguration(&configuration)
	}

}
func SaveConfiguration(c *LocalConfiguration) {
	log.Logger.Infof("m=SaveConfiguration, status=begin, configuration=%+v", c)
	if len(c.Envs) == 0 {
		c.Envs = NewEmptyEnv()
	}

	js,_ := json.Marshal(&c)
	log.Logger.Infof("m=SaveConfiguration, status=save, data=%s", js)

	f, err := os.OpenFile(utils.GetPath(confPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	defer func(){
		f.Close()
	}()
	if err != nil {
		log.Logger.Errorf("status=error-to-create-conf-file, err=%v", err)
		return
	}
	wr := bufio.NewWriter(f)
	defer func(){
		wr.Flush()
	}()
	enc := json.NewEncoder(wr)
	err = enc.Encode(c)
	if err != nil {
		log.Logger.Errorf("status=error-to-encode, error=%v", err)
	}

	log.Logger.Infof("m=SaveConfiguration, status=success")

}

func GetConfiguration() LocalConfiguration {
	LoadConfiguration()
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

func (lc *LocalConfiguration) AddHostnameToEnv(env string, hostname *HostnameVo){
	log.Logger.Infof("m=AddHostnameToEnv, status=begin, env=%+v, hostname=%+v", env, hostname)
	foundEnv := lc.GetEnv(env)
	t := append(foundEnv.Hostnames, *hostname)
	(*foundEnv).Hostnames = t
	foundEnv.Name = "tmp"
	log.Logger.Infof("m=AddHostnameToEnv, status=success, lc=%+v, foundEnv=%+v, hostnames=%+v", lc, foundEnv, lc.Envs[0].Hostnames)
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

func (lc *LocalConfiguration) AddEnv(env EnvVo){
	lc.Envs = append(lc.Envs, env)
	SaveConfiguration(lc)
}

func (lc *LocalConfiguration) RemoveEnv(index int){
	lc.Envs = append(lc.Envs[:index], lc.Envs[index+1:]...)
	SaveConfiguration(lc)
}

func (lc *LocalConfiguration) AddDns(dns [4]byte){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers, dns)
	SaveConfiguration(lc)
}

func (lc *LocalConfiguration) RemoveDns(index int){
	lc.RemoteDnsServers = append(lc.RemoteDnsServers[:index], lc.RemoteDnsServers[index+1:]...)
	SaveConfiguration(lc)
}


func AddHostname(envName string, hostname HostnameVo){
	log.Logger.Infof("m=AddHostname, status=begin, evnName=%s, hostname=%+v", envName, hostname)
	configuration.AddHostnameToEnv(envName, &hostname)
	SaveConfiguration(&configuration)
	log.Logger.Infof("m=AddHostname, status=success, configuration=%+v", configuration)
}

func RemoveHostname(envIndex int, hostIndex int){
	env := configuration.Envs[envIndex];
	t := append(env.Hostnames[:hostIndex], env.Hostnames[hostIndex+1:]...)
	env.Hostnames = t
	SaveConfiguration(&configuration)
}

func NewEmptyEnv() []EnvVo {
	return []EnvVo{{Hostnames:[]HostnameVo{}, Name:""}}
}