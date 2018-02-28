package main

import (
	"log"
	"net/http"

	"github.com/romosborne/nozzy-tasks/models"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName = "./tasks.db"
)

func main() {

	db, err := models.NewDB(databaseName)
	if err != nil {
		log.Panic(err)
	}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
