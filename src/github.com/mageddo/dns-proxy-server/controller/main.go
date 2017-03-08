package controller

import (
	"net/http"
)

func MapRequests() {
	// this is a placebo to execute inits from this package
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, r.URL.Path)
	}
}