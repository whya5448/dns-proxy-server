package hostname

import (
	"net/http"
	"encoding/json"
	"github.com/mageddo/dns-proxy-server/events/local"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
	r "github.com/mageddo/dns-proxy-server/controller"
)

func init(){
	r.Get("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		json.NewEncoder(res).Encode(local.GetConfiguration(ctx))
	})

	r.Post("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=create-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, host=%+v", hostname)
		err := local.GetConfiguration(ctx).AddHostname(ctx, hostname.Env, hostname)
		if err != nil {
			logger.Infof("m=/hostname/, status=error, action=create-hostname, err=%+v", err)
			r.BadRequest(res, err.Error())
		}
		logger.Infof("m=/hostname/, status=success, action=create-hostname")
	})

	r.Delete("/hostname/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		logger := log.GetLogger(ctx)
		res.Header().Add("Content-Type", "application/json")
		logger.Infof("m=/hostname/, status=begin, action=delete-hostname")
		var hostname local.HostnameVo
		json.NewDecoder(req.Body).Decode(&hostname)
		logger.Infof("m=/hostname/, status=parsed-host, action=delete-hostname, host=%+v", hostname)
		err := local.GetConfiguration(ctx).RemoveHostnameByEnvAndHostname(ctx, hostname.Env, hostname.Hostname)
		if err != nil {
			logger.Infof("m=/hostname/, status=error, action=delete-hostname, err=%+v", err)
			r.BadRequest(res, err.Error())
		}
		logger.Infof("m=/hostname/, status=success, action=delete-hostname")
	})
}