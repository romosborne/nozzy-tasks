package main

import (
	"net/http"
	"nozzy-tasks/models"
)

type Route struct {
	Name             string
	Method           string
	Pattern          string
	ApiAuthenticated bool
	WebAuthenticated bool
	HandlerFunc      func(*models.Env) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"GET",
		"/login",
		false,
		false,
		WebLogin,
	},
	Route{
		"Auth",
		"GET",
		"/auth",
		false,
		false,
		WebAuth,
	},
	Route{
		"Secure",
		"GET",
		"/secure",
		false,
		true,
		WebSecure,
	},
	Route{
		"Index",
		"GET",
		"/",
		false,
		false,
		Index,
	},
	Route{
		"TaskIndex",
		"GET",
		"/tasks",
		true,
		false,
		TaskIndex,
	},
	Route{
		"TaskShow",
		"GET",
		"/tasks/{taskId}",
		true,
		false,
		TaskShow,
	},
	Route{
		"TaskCreate",
		"POST",
		"/tasks",
		true,
		false,
		TaskCreate,
	},
}
