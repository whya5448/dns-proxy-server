package controller

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/dns-proxy-server/events/local"
)

func init(){
	Get("/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Location", "/static")
		res.WriteHeader(301)
	})

	Get("/configuration/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Header().Add("Content-Type", "application/json")
		utils.GetJsonEncoder(res).Encode(local.GetConfiguration(ctx))
	})
}