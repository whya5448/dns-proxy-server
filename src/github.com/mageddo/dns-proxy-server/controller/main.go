package controller

import (
	"net/http"
	"github.com/mageddo/dns-proxy-server/utils"
	"github.com/mageddo/log"
)

func Map() {
	http.HandleFunc("/hello", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	}))
	http.HandleFunc("/static", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		staticPath := utils.GetPath("/static")
		log.Logger.Infof("path=%v", staticPath)
		hd := http.FileServer(http.Dir(staticPath))
		hd.ServeHTTP(res, req)
	}))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//m := validPath.FindStringSubmatch(r.URL.Path)
		//if m == nil {
		//	http.NotFound(w, r)
		//	return
		//}
		//fn(w, r, m[2])
		fn(w, r, r.URL.Path)
	}
}