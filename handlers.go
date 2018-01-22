package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		id, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		task, _ := env.db.SingleTask(id)
		json.NewEncoder(w).Encode(task)
	}
}

func TaskCreate(env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		if err := r.Body.Close(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		if err := json.Unmarshal(body, &task); err != nil {
			w.Header().Set("Content-type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err = env.db.CreateTask(&task)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(task); err != nil {
			panic(err)
		}
	}
}
