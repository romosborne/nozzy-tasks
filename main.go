package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	env := &models.Env{}

	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &env)
	if err != nil {
		fmt.Printf("File parse error: %v\n", err)
		os.Exit(1)
	}

	env.Db = db
	env.SessionKey = []byte(RandToken(64))

	router := NewRouter(env)

	log.Fatal(http.ListenAndServe(":8080", router))
}
