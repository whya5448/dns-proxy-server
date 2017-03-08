package controller

import "net/http"

func init(){
	http.HandleFunc("/hello/", makeHandler(func(res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	}))
}