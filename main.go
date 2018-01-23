package main

import (
	"log"
	"net/http"

	"nozzy-tasks/models"

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

	env := &models.Env{
		Db:         db,
		SessionKey: []byte(RandToken(64)),
	}

	router := NewRouter(env)

	log.Fatal(http.ListenAndServe(":8080", router))
}
