package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/events/local"
	"encoding/json"
)

func init(){
	http.HandleFunc("/hello/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	}))

	http.HandleFunc("/hostname/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		switch req.Method {
		case "GET":
			json.NewEncoder(res).Encode(local.GetConfiguration())
			return
			break
		case "POST":
			var hostname local.HostnameVo
			json.NewDecoder(req.Body).Decode(&hostname)

			conf := local.GetConfiguration()
			local.AddHostname(*conf.GetActiveEnv(), hostname)
			return
			break
		}

		res.Write([]byte("It works from controller!!!"))
	}))
}