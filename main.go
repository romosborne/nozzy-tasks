package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/romosborne/nozzy-tasks/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	bindAddress string
)

const (
	databaseName = "./tasks.db"
)

func init() {
	flag.StringVar(&bindAddress, "bind", ":8082", "Specify the ip and port to listen on")
}

func main() {
	flag.Parse()

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

	log.Println("Binding on", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, router))
}
