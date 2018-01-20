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

func main() {

	db, err := models.NewDB(databaseName)
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db: db}

	http.HandleFunc("/", env.allTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (env *Env) allTasks(w http.ResponseWriter, r *http.Request) {
	tasks, _ := env.db.AllTasks()
	json.NewEncoder(w).Encode(tasks)
}
