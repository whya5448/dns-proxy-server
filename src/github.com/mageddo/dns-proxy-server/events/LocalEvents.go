package events

import (
	"github.com/mageddo/dns-proxy-server/utils"
	"encoding/json"
	"os"
	"github.com/mageddo/log"
	"bufio"
)

var configuration = utils.LocalConfiguration{
	Envs: *new([]utils.EnvVo),
	RemoteDnsServers: *new([][4]byte),
}
//var cache = make(map[string]utils.HostnameVo)


func x(){


}

func LoadConfiguration(){

	//configuration.RemoteDnsServers = append(configuration.RemoteDnsServers, [4]byte{1,2,3,4})

	var confPath string = "/app/conf/config.json"
	if _, err := os.Stat(confPath); err == nil {

		f, _ := os.Open(confPath)
		dec := json.NewDecoder(f)
		dec.Decode(&configuration)

	}else{
		err := os.MkdirAll("/app/conf", 0755)
		if err != nil {
			log.Logger.Errorf("status=error-to-create-conf-folder, err=%v", err)
			return
		}
		f, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Logger.Errorf("status=error-to-create-conf-file, err=%v", err)
			return
		}
		wr := bufio.NewWriter(f)
		enc := json.NewEncoder(wr)
		err = enc.Encode(configuration)
		log.Logger.Errorf("%v", configuration)
		if err != nil {
			log.Logger.Errorf("status=error-to-encode, error=%v", err)
		}
		wr.Flush()
		f.Close()
	}


}
