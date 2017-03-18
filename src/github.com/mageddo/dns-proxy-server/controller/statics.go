package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
)

func init(){

	http.HandleFunc("/static/", func(res http.ResponseWriter, req *http.Request){
		logger := log.GetLogger(log.GetContext())

		staticPath := utils.GetPath("/")
		logger.Infof("urlPath=%s", req.URL.Path)
		hd := http.FileServer(http.Dir(staticPath))
		hd.ServeHTTP(res, req)
	})
}
