package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nozzy-tasks/models"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type JwtToken struct {
	Token string `json:"token"`
}

func Authenticate(env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		pass, err := env.db.CheckPassword(user.Username, user.Password)
		if err != nil {
			fmt.Println(err)
		}

		if pass != true {
			fmt.Fprint(w, "Invalid username or password")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"password": user.Password,
		})
		tokenString, error := token.SignedString([]byte("secret"))
		if error != nil {
			fmt.Println(error)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	}
}

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
