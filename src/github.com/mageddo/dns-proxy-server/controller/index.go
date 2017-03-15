package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/events/local"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
)

func init(){
	http.HandleFunc("/hello/", makeHandler(func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	}))

	http.HandleFunc("/hostname/", makeHandler(func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		switch req.Method {
		case "GET":
			json.NewEncoder(res).Encode(local.GetConfiguration())
			return
			break
		case "POST":
			log.GetLogger(ctx).Infof("post")
			var hostname local.HostnameVo
			json.NewDecoder(req.Body).Decode(&hostname)
			conf := local.GetConfiguration()
			local.AddHostname((*conf.GetActiveEnv()).Name, hostname)
			return
			break
		}

	}))


	http.HandleFunc("/hostname/new/", makeHandler(func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/new/, status=begin")
		switch req.Method {
		case "POST":
			var hostname local.HostnameVo
			json.NewDecoder(req.Body).Decode(&hostname)
			conf := local.GetConfiguration()
			env := (*conf.GetActiveEnv()).Name
			logger.Infof("m=/hostname/new/, status=parsed-host, env=%+v, host=%+v", env, hostname)
			local.AddHostname(env, hostname)
		}
		log.GetLogger(ctx).Infof("m=/hostname/new/, status=success")
	}))
}