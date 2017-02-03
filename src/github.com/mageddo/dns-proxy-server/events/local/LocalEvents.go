package local

import (
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
)

var confPath string = "conf/config.json"
var configuration = LocalConfiguration{
	Envs: make([]EnvVo, 0),
	RemoteDnsServers: make([][4]byte, 0),
}

func LoadConfiguration(){

	if _, err := os.Stat(confPath); err == nil {

		f, _ := os.Open(confPath)

		defer func(){
			f.Close()
		}()

		dec := json.NewDecoder(f)
		dec.Decode(&configuration)

	}else{
		err := os.MkdirAll("conf", 0755)
		if err != nil {
			log.Logger.Errorf("status=error-to-create-conf-folder, err=%v", err)
			return
		}
		SaveConfiguration(&configuration)
	}

}
func SaveConfiguration(configuration *LocalConfiguration) {

	f, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
	err = enc.Encode(configuration)
	if err != nil {
		log.Logger.Errorf("status=error-to-encode, error=%v", err)
	}

}

func GetConfiguration() LocalConfiguration {
	return configuration
}


type LocalConfiguration struct {
	RemoteDnsServers [][4]byte
	Envs []EnvVo
	ActiveEnv string
}

type EnvVo struct {
	Name string
	Hostnames []HostnameVo
}

type HostnameVo struct {
	Hostname string
	Ip [4]byte
	Ttl int
}

func (lc *LocalConfiguration) GetEnv(envName string) (*EnvVo) {

	for _, env := range lc.Envs {
		if env.Name == envName {
			return &env
		}
	}
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

func (lc *LocalConfiguration) AddEnv(env EnvVo){
	configuration.Envs = append(configuration.Envs, env)
	SaveConfiguration(&configuration)
}

func (lc *LocalConfiguration) RemoveEnv(index int){
	append(configuration.Envs[:index], configuration.Envs[index+1:]...)
	SaveConfiguration(&configuration)
}

func (lc *LocalConfiguration) AddDns(dns [4]byte){
	configuration.RemoteDnsServers = append(configuration.RemoteDnsServers, dns)
	SaveConfiguration(&configuration)
}

func (lc *LocalConfiguration) RemoveDns(index int){
	append(configuration.RemoteDnsServers[:index], configuration.RemoteDnsServers[index+1:]...)
	SaveConfiguration(&configuration)
}


func AddHostname(env EnvVo, hostname HostnameVo){
	foundEnv := configuration.GetEnv(env.Name)
	env.Hostnames = append(foundEnv.Hostnames, hostname)
	SaveConfiguration(&configuration)
}

func RemoveHostname(envIndex int, hostIndex int){
	env := configuration.Envs[envIndex];
	append(env.Hostnames[:hostIndex], env.Hostnames[hostIndex+1:]...)
	SaveConfiguration(&configuration)
}