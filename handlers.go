package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/romosborne/nozzy-tasks/models"

	"github.com/futurenda/google-auth-id-token-verifier"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	store *sessions.CookieStore
)

type viewBag struct {
	Link  string
	Email string
}

// RandToken returns a random string of length l
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Web returns the web interface
func Web(_ *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("./static/web.html")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(content)
	}
}

// ApiAuth takes a google jwt, validates it, and returns a authtoken for the user
func ApiAuth(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
			return
		}
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			jwt := bearerToken[1]

			v := googleAuthIDTokenVerifier.Verifier{}
			err := v.VerifyIDToken(jwt, []string{
				env.OauthClientID,
			})
			if err != nil {
				json.NewEncoder(w).Encode(models.Exception{Message: "Invalid authorization token"})
				return
			}

			claimSet, err := googleAuthIDTokenVerifier.Decode(jwt)

			authToken := RandToken(64)

			env.Db.CreateOrSetAuthToken(claimSet.Sub, claimSet.Email, authToken)

			fmt.Fprintf(w, "%s", authToken)
		}
	}
}

// TaskIndex returns all projects and tasks
func TaskIndex(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=UTF-8")
		user := getUser(env, r)
		tasks, _ := env.Db.AllTasks(user.ID)
		json.NewEncoder(w).Encode(tasks)
	}
}

// TaskShow returns a single task
func TaskShow(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=UTF-8")

		vars := mux.Vars(r)
		taskID := vars["taskId"]
		id, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		userID := getUser(env, r).ID

		task, _ := env.Db.SingleTask(id, userID)
		json.NewEncoder(w).Encode(task)
	}
}

// TaskCreate creates a task
func TaskCreate(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskRequest models.TaskRequest
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))

		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			fmt.Fprint(w, err)
			return
		}

		if err := json.Unmarshal(body, &taskRequest); err != nil {
			w.Header().Set("Content-type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		if taskRequest.NewProjectName != "" {
			project := models.Project{
				Name:   taskRequest.NewProjectName,
				UserID: getUser(env, r).ID,
			}

			err = env.Db.CreateProject(&project)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}

			taskRequest.ProjectID = project.ID
		}

		task := models.Task{
			Title:     taskRequest.Title,
			ProjectID: taskRequest.ProjectID,
		}

		err = env.Db.CreateTask(&task)
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

func TaskComplete(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tcr models.TaskCompletionRequest
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))

		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			fmt.Fprint(w, err)
			return
		}

		if err := json.Unmarshal(body, &tcr); err != nil {
			w.Header().Set("Content-type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err = env.Db.SetTaskCompletion(getUser(env, r).ID, tcr.TaskID, tcr.Completed)

		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}
}

// ProjectCreate creates a project
func ProjectCreate(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var project models.Project
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		if err := r.Body.Close(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		if err := json.Unmarshal(body, &project); err != nil {
			w.Header().Set("Content-type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		project.UserID = getUser(env, r).ID

		err = env.Db.CreateProject(&project)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(project); err != nil {
			panic(err)
		}
	}
}

func getUser(env *models.Env, r *http.Request) *models.User {
	return r.Context().Value(env.ContextKey).(*models.User)
}
