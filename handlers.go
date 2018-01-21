package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(_ *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome!")
	}
}

func TaskIndex(env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, _ := env.db.AllTasks()
		json.NewEncoder(w).Encode(tasks)
	}
}

func TaskShow(env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["taskId"]
		id, err := strconv.Atoi(taskID)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		task, _ := env.db.SingleTask(id)
		json.NewEncoder(w).Encode(task)
	}
}
