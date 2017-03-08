package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
)

func MapRequests() {
	http.HandleFunc("/hello/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	}))
	http.HandleFunc("/static/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		staticPath := utils.GetPath("/")
		log.Logger.Infof("path=%v", staticPath)
		hd := http.FileServer(http.Dir(staticPath))
		hd.ServeHTTP(res, req)
	}))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, r.URL.Path)
	}
}