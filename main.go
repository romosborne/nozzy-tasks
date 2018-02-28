package main

import (
<<<<<<< HEAD
=======
	"flag"
>>>>>>> master
	"log"
	"net/http"

	"github.com/romosborne/nozzy-tasks/models"

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

<<<<<<< HEAD
	router := NewRouter()
=======
func main() {
	flag.Parse()

	if oauthClientID == "" {
		log.Println("Please specify your oauth2 client id")
		os.Exit(1)
	}

	db, err := models.NewDB(databaseName)
	if err != nil {
		log.Panic(err)
	}

	env := &models.Env{
		OauthClientID: oauthClientID,
		Db:            db,
		SessionKey:    []byte(RandToken(64)),
	}

	router := NewRouter(env)
>>>>>>> master

	log.Println("Binding on", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, router))
}
