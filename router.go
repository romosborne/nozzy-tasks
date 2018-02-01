package main

import (
	"net/http"
	"nozzy-tasks/middleware"
	"nozzy-tasks/models"
	"strings"

	"github.com/gorilla/mux"
)

// NewRouter creates a router using the routes defined in Routes
func NewRouter(env *models.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(env)

		if strings.Contains(route.Pattern, "/api/") {
			handler = middleware.APIValidate(env, handler)
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
