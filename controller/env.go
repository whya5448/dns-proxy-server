package controller

import (
	log "github.com/mageddo/go-logging"
	"net/http"
	"encoding/json"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/utils"
)

func init(){

	Get("/env/active", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			utils.GetJsonEncoder(res).Encode(local.EnvVo{Name: conf.ActiveEnv})
			return
		}
		confLoadError(res)
	})

	Put("/env/active", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){

		logger := log.NewLog(ctx)
		logger.Infof("m=/env/active/, status=begin")
		res.Header().Add("Content-Type", "application/json")

		var envVo local.EnvVo
		json.NewDecoder(req.Body).Decode(&envVo)
		logger.Infof("m=/env/active/, status=parsed-host, env=%+v", envVo)

		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.SetActiveEnv(ctx, envVo); err != nil {
				logger.Infof("m=/env/, status=error, action=create-env, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/env/active/, status=success, action=active-env")
			return
		}
		confLoadError(res)

	})

	Get("/env/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			utils.GetJsonEncoder(res).Encode(conf.Envs)
			return
		}
		confLoadError(res)
	})

	Post("/env/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.NewLog(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/env/, status=begin, action=create-env")
		var envVo local.EnvVo
		json.NewDecoder(req.Body).Decode(&envVo)
		logger.Infof("m=/env/, status=parsed-host, env=%+v", envVo)
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.AddEnv(ctx, envVo);  err != nil {
				logger.Infof("m=/env/, status=error, action=create-env, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/env/, status=success, action=create-env")
			return
		}
		confLoadError(res)
	})

	Delete("/env/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.NewLog(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/env/, status=begin, action=delete-env")
		var env local.EnvVo
		json.NewDecoder(req.Body).Decode(&env)
		logger.Infof("m=/env/, status=parsed-host, action=delete-env, env=%+v", env)
		if conf, _ := local.LoadConfiguration(ctx); conf != nil {
			if err := conf.RemoveEnvByName(ctx, env.Name);  err != nil {
				logger.Infof("m=/env/, status=error, action=delete-env, err=%+v", err)
				BadRequest(res, err.Error())
				return
			}
			logger.Infof("m=/env/, status=success, action=delete-env")
			return
		}
		confLoadError(res)
	})
}

func confLoadError(res http.ResponseWriter){
	BadRequest(res, "Could not load conf")
}
