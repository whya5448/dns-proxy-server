package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
)

func init(){
	http.HandleFunc("/static/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		staticPath := utils.GetPath("/")
		log.Logger.Infof("path=%v", staticPath)
		hd := http.FileServer(http.Dir(staticPath))
		hd.ServeHTTP(res, req)
	}))
}
