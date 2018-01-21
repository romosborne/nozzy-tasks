package main

import (
	"net/http"
)

type Route struct {
	Name          string
	Method        string
	Pattern       string
	Authenticated bool
	HandlerFunc   func(*Env) http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Authenticate",
		"POST",
		"/authenticate",
		false,
		Authenticate,
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
		"/tasks",
		true,
		TaskIndex,
	},
	Route{
		"TaskShow",
		"GET",
		"/tasks/{taskId}",
		true,
		TaskShow,
	},
}
