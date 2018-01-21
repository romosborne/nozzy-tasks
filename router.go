package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(env *Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = Logger(route.HandlerFunc(env), route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
