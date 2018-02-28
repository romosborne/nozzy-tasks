package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/romosborne/nozzy-tasks/server"
	"github.com/romosborne/nozzy-tasks/services"

	_ "github.com/mattn/go-sqlite3"
)

var (
	bindAddress   string
	oauthClientID string
)

const (
	databaseName = "./tasks.db"
)

func init() {
	flag.StringVar(&bindAddress, "bind", ":8082", "Specify the ip and port to listen on")
	flag.StringVar(&oauthClientID, "oauth-id", "", "Specify your oauth2 client id")
}

func main() {
	flag.Parse()

	if oauthClientID == "" {
		log.Println("Please specify your oauth2 client id")
		os.Exit(1)
	}

	sql, err := services.NewSQL(databaseName)
	if err != nil {
		log.Panic(err)
	}

	server := server.NewRouter(sql, oauthClientID)

	log.Println("Binding on", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, server))
}
