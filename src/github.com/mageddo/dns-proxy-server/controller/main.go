package controller

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
	"encoding/json"
)


type Method string
const (
POST Method = "POST"
GET Method = "GET"
PUT Method = "PUT"
PATCH Method = "PATCH"
DELETE Method = "DELETE"
)

type Map struct {
	method Method
	path string
}

type Message struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func BadRequest(w http.ResponseWriter, msg string){
	RespMessage(w, 400, msg)
}

func RespMessage(w http.ResponseWriter, code int, msg string){
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Message{Code:code, Message:msg})
}

var maps = make(map[string]map[string]func(context.Context, http.ResponseWriter, *http.Request, string))

func MapRequests() {
	// this is a placebo to execute inits from this package
}

func Post(path string, fn func(context.Context, http.ResponseWriter, *http.Request, string)) {
	MapReq(POST, path, fn)
}

func Get(path string, fn func(context.Context, http.ResponseWriter, *http.Request, string)) {
	MapReq(GET, path, fn)
}

func Put(path string, fn func(context.Context, http.ResponseWriter, *http.Request, string)) {
	MapReq(PUT, path, fn)
}

func MapReq(method Method, path string, fn func(context.Context, http.ResponseWriter, *http.Request, string)) {

	_, mapped := maps[path]
	if !mapped {

		maps[path] = make(map[string]func(context.Context, http.ResponseWriter, *http.Request, string))

		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			_, matched := maps[urlPath][r.Method]
			ctx := log.GetContext()
			logger := log.GetLogger(ctx)

			logger.Debugf("m=MapReq, status=begin, matched=%t, url=%s, method=%s", matched, urlPath,  r.Method)
			if matched {
				function := maps[urlPath][r.Method]
				function(ctx, w, r, urlPath)
				logger.Debugf("m=MapReq, status=success, url=%s %s", r.Method, urlPath)
			}else{
				logger.Debugf("m=MapReq, status=not-found, url=%s %s", r.Method, urlPath)
				http.NotFound(w, r)
			}

		})
	}
	log.Logger.Debugf("m=MapReq, status=mapping, url=%s %s", method, path)
	maps[path][string(method)] = fn

}