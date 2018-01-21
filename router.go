package main

import (
	"net/http"
	"nozzy-tasks/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(env *Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(env)

		if route.Authenticated {
			handler = middleware.Validate(handler)
		}

		handler = middleware.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
