package main

import (
	"net/http"
	"nozzy-tasks/models"
)

// A Route is a http route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*models.Env) http.HandlerFunc
}

// Routes is a collection os routes
type Routes []Route

var routes = Routes{
	Route{
		"Web",
		"GET",
		"/",
		Web,
	},
	Route{
		"ApiAuth",
		"GET",
		"/get_token",
		ApiAuth,
	},
	Route{
		"TaskIndex",
		"GET",
		"/api/tasks",
		TaskIndex,
	},
	Route{
		"TaskShow",
		"GET",
		"/api/tasks/{taskId}",
		TaskShow,
	},
	Route{
		"TaskCreate",
		"POST",
		"/api/tasks",
		TaskCreate,
	},
	Route{
		"TaskCompletion",
		"POST",
		"/api/tasks/completion",
		TaskComplete,
	},
	Route{
		"ProjectCreate",
		"POST",
		"/api/project",
		ProjectCreate,
	},
}
