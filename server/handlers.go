package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/romosborne/nozzy-tasks/middleware"
	"github.com/romosborne/nozzy-tasks/models"
	"github.com/romosborne/nozzy-tasks/services"

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
func Web(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/web.html")
	if err != nil {
		fmt.Println(err)
	}

	clientID := getOauthClientID(r)

	t.Execute(w, clientID)
}

// APIAuth takes a google jwt, validates it, and returns a authtoken for the user
func APIAuth(w http.ResponseWriter, r *http.Request) {
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
			getOauthClientID(r),
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.Exception{Message: "Invalid authorization token"})
			return
		}

		claimSet, err := googleAuthIDTokenVerifier.Decode(jwt)

		authToken := RandToken(64)

		sqlService := getSQLService(r)
		sqlService.CreateOrSetAuthToken(claimSet.Sub, claimSet.Email, authToken)

		fmt.Fprintf(w, "%s", authToken)
	}
}

// TaskIndex returns all projects and tasks
func TaskIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	user := getUser(r)
	sqlService := getSQLService(r)
	tasks, _ := sqlService.AllTasks(user.ID)
	json.NewEncoder(w).Encode(tasks)
}

// TaskShow returns a single task
func TaskShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	taskID := vars["taskId"]
	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	userID := getUser(r).ID

	sqlService := getSQLService(r)
	task, _ := sqlService.SingleTask(id, userID)
	json.NewEncoder(w).Encode(task)
}

// TaskDelete handles delete task requests
func TaskDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	taskID := vars["taskId"]
	id, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	userID := getUser(r).ID

	sqlService := getSQLService(r)
	err = sqlService.DeleteTask(id, userID)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// TaskCreate creates a task
func TaskCreate(w http.ResponseWriter, r *http.Request) {
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

	sqlService := getSQLService(r)

	if taskRequest.NewProjectName != "" {
		project := models.Project{
			Name:   taskRequest.NewProjectName,
			UserID: getUser(r).ID,
		}

		err = sqlService.CreateProject(&project)
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

	err = sqlService.CreateTask(&task)
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

// TaskComplete handles task completion requests
func TaskComplete(w http.ResponseWriter, r *http.Request) {
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

	sqlService := getSQLService(r)
	err = sqlService.SetTaskCompletion(getUser(r).ID, tcr.TaskID, tcr.Completed)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

// ProjectCreate creates a project
func ProjectCreate(w http.ResponseWriter, r *http.Request) {
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

	project.UserID = getUser(r).ID
	sqlService := getSQLService(r)

	err = sqlService.CreateProject(&project)
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

func ProjectDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	projectID := vars["projectId"]
	id, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	userID := getUser(r).ID

	sqlService := getSQLService(r)
	err = sqlService.DeleteProject(id, userID)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getUser(r *http.Request) *models.User {
	return r.Context().Value(middleware.UserContextKey).(*models.User)
}

func getSQLService(r *http.Request) *services.SQL {
	return r.Context().Value(middleware.SQLServiceKey).(*services.SQL)
}

func getOauthClientID(r *http.Request) string {
	return r.Context().Value(middleware.OauthClientIDKey).(string)
}
