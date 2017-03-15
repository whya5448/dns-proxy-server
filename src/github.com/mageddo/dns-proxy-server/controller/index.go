package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/events/local"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
)

func init(){
	Get("/hello/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	})

	Get("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		json.NewEncoder(res).Encode(local.GetConfiguration())
	})

	Post("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/new/, status=begin")
		switch req.Method {
		case "POST":
			var hostname local.HostnameVo
			json.NewDecoder(req.Body).Decode(&hostname)
			logger.Infof("m=/hostname/new/, status=parsed-host, host=%+v", hostname)
			local.AddHostname(hostname.Env, hostname)
		}
		log.GetLogger(ctx).Infof("m=/hostname/new/, status=success")
	})
}