package main

import (
	"net/http"
	"nozzy-tasks/models"
)

type Route struct {
	Name             string
	Method           string
	Pattern          string
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
		WebLogin,
	},
	Route{
		"Auth",
		"GET",
		"/auth",
		false,
		WebAuth,
	},
	Route{
		"Secure",
		"GET",
		"/secure",
		true,
		WebSecure,
	},
	Route{
		"Index",
		"GET",
		"/",
		false,
		Index,
	},
	Route{
		"TaskIndex",
		"GET",
		"/api/tasks",
		false,
		TaskIndex,
	},
	Route{
		"TaskShow",
		"GET",
		"/api/tasks/{taskId}",
		false,
		TaskShow,
	},
	Route{
		"TaskCreate",
		"POST",
		"/api/tasks",
		false,
		TaskCreate,
	},
	Route{
		"ProjectCreate",
		"POST",
		"/api/project",
		false,
		ProjectCreate,
	},
}
