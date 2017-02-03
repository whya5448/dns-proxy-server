package local

import (
	"github.com/mageddo/dns-proxy-server/utils"
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
)

var confPath string = "conf/config.json"
var configuration = utils.LocalConfiguration{
	Envs: make([]utils.EnvVo, 0),
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
		saveConfiguration(&configuration)
	}

}
func saveConfiguration(configuration *utils.LocalConfiguration) {

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

func GetConfiguration() utils.LocalConfiguration {
	return configuration
}