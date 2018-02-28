package server

import (
	"net/http"
	"strings"

	"github.com/romosborne/nozzy-tasks/services"

	"github.com/romosborne/nozzy-tasks/middleware"

	"github.com/gorilla/mux"
)

// NewRouter creates a router using the routes defined in Routes
func NewRouter(sql *services.SQL, oauthClientID string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if strings.Contains(route.Pattern, "/api/") {
			handler = middleware.APIValidate(handler)
		}

		handler = middleware.Logger(handler, route.Name)
		handler = middleware.AddContext(handler, sql, oauthClientID)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
