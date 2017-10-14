package controller

import (
	"net/http"
	"encoding/json"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	log "github.com/mageddo/go-logging"
	"fmt"
)

func init(){
	Get("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			envName := req.URL.Query().Get("env")
			if env, _ := conf.GetEnv(envName);  env != nil {
				json.NewEncoder(res).Encode(env)
				return
			}
			BadRequest(res, fmt.Sprintf("Env %s not found", envName))
			return
		}
		confLoadError(res)
	})

	Get("/hostname/find/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			env := req.URL.Query().Get("env")
			hostname := req.URL.Query().Get("hostname")
			var err error
			var hostnames *[]local.HostnameVo
			if hostnames, err = conf.FindHostnameByNameAndEnv(ctx, env, hostname);  err == nil {
				json.NewEncoder(res).Encode(hostnames)
				return
			}
			BadRequest(res, fmt.Sprintf(err.Error()))
			return
		}
		confLoadError(res)
	})

	Post("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.NewLog(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=create-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, host=%+v", hostname)
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.AddHostname(ctx, hostname.Env, hostname); err != nil {
				logger.Infof("m=/hostname/, status=error, action=create-hostname, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/hostname/, status=success, action=create-hostname")
			return
		}
		confLoadError(res)
	})

	Put("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.NewLog(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=update-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, host=%+v", hostname)
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.UpdateHostname(ctx, hostname.Env, hostname);  err != nil {
				logger.Infof("m=/hostname/, status=error, action=update-hostname, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/hostname/, status=success, action=update-hostname")
			return
		}
		confLoadError(res)
	})

	Delete("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.NewLog(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=delete-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, action=delete-hostname, host=%+v", hostname)
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.RemoveHostnameByEnvAndHostname(ctx, hostname.Env, hostname.Hostname);  err != nil {
				logger.Infof("m=/hostname/, status=error, action=delete-hostname, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/hostname/, status=success, action=delete-hostname")
			return
		}
		confLoadError(res)
	})
}
