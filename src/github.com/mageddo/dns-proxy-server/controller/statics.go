package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
	"golang.org/x/net/context"
)

func init(){
	Get("/static/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		staticPath := utils.GetPath("/")
		log.Logger.Infof("path=%v", staticPath)
		hd := http.FileServer(http.Dir(staticPath))
		hd.ServeHTTP(res, req)
	})
}
