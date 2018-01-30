package main

import (
	"net/http"
	"nozzy-tasks/models"
)

// A Route is a http route
type Route struct {
	Name             string
	Method           string
	Pattern          string
	WebAuthenticated bool
	HandlerFunc      func(*models.Env) http.HandlerFunc
}

// Routes is a collection os routes
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
		"Tasks",
		"GET",
		"/tasks",
		true,
		WebTasks,
	},
	Route{
		"NewProject",
		"GET",
		"/newProject",
		true,
		WebNewProject,
	},
	Route{
		"NewProjectPost",
		"POST",
		"/newProject",
		true,
		WebNewProjectPost,
	},
	Route{
		"NewTask",
		"GET",
		"/newTask",
		true,
		WebNewTask,
	},
	Route{
		"NewTaskPost",
		"POST",
		"/newTask",
		true,
		WebNewTaskPost,
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
