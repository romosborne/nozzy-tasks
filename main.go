package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/romosborne/nozzy-tasks/models"
	"github.com/romosborne/nozzy-tasks/server"
	"github.com/romosborne/nozzy-tasks/services"

	_ "github.com/mattn/go-sqlite3"
)

var (
	bindAddress   string
	oauthClientID string
)

const (
	databaseName = "/config/tasks.db"
)

func init() {
	flag.StringVar(&bindAddress, "bind", ":8082", "Specify the ip and port to listen on")
	flag.StringVar(&oauthClientID, "oauth-id", "", "Specify your oauth2 client id")
}

func main() {
	var config models.Config

	configFile, err := ioutil.ReadFile("/config/config.json")
	if err != nil {
		log.Println("Error opening config file")
		os.Exit(1)
	}

	_ = json.Unmarshal([]byte(configFile), &config)

	sql, err := services.NewSQL(databaseName)
	if err != nil {
		log.Panic(err)
	}

	server := server.NewRouter(sql, config.Cid)

	log.Println("Binding on", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, server))
}
