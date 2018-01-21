package main

import (
	"encoding/json"
	"log"
	"net/http"

	"nozzy-tasks/models"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName = "./tasks.db"
)

type Env struct {
	db models.Datastore
}

func (env *Env) allTasks(w http.ResponseWriter, r *http.Request) {
	tasks, _ := env.db.AllTasks()
	json.NewEncoder(w).Encode(tasks)
}

func main() {

	db, err := models.NewDB(databaseName)
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db: db}

	router := NewRouter(env)

	log.Fatal(http.ListenAndServe(":8080", router))
}
