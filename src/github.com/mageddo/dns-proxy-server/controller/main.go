package controller

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/mageddo/log"
)

func MapRequests() {
	// this is a placebo to execute inits from this package
}

func makeHandler(fn func(context.Context, http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(log.GetContext(), w, r, r.URL.Path)
	}
}