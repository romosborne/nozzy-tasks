package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*Env) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TaskIndex",
		"GET",
		"/tasks",
		TaskIndex,
	},
	Route{
		"TaskShow",
		"GET",
		"/tasks/{taskId}",
		TaskShow,
	},
}
