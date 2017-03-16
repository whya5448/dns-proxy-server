package controller

import (
	"net/http"
	"golang.org/x/net/context"
)

func init(){
	Get("/hello/", func(ctx context.Context, res http.ResponseWriter, req *http.Request, url string){
		res.Write([]byte("It works from controller!!!"))
	})
}