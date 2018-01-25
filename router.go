package main

import (
	"net/http"
	"nozzy-tasks/middleware"
	"nozzy-tasks/models"
	"strings"

	"github.com/gorilla/mux"
)

func NewRouter(env *models.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(env)

		if strings.Contains(route.Pattern, "/api/") {
			handler = middleware.ApiValidate(env, handler)
		}

		if route.WebAuthenticated {
			handler = middleware.WebValidate(env, handler)
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
