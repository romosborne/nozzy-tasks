package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"nozzy-tasks/models"
	"strconv"
	"strings"

	"github.com/futurenda/google-auth-id-token-verifier"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func createOauthConf(env *models.Env) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     env.OauthClientID,
		ClientSecret: env.OauthClientSecret,
		RedirectURL:  env.OauthRedirectURL,
		Scopes: []string{
			"openid",
			"email",
		},
		Endpoint: google.Endpoint,
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

			env.Db.SetAuthToken(claimSet.Sub, claimSet.Email, authToken)

			fmt.Fprintf(w, "%s", authToken)
		}
	}
}

// WebAuth is the login call for the web interface
func WebAuth(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := sessions.NewCookieStore(env.SessionKey)
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		retrievedState := session.Values["state"].(string)
		queryState := r.URL.Query()["state"][0]
		if retrievedState != queryState {
			log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
			http.Error(w, "Invalid session state", http.StatusUnauthorized)
			return
		}

		conf := createOauthConf(env)

		code := r.URL.Query()["code"][0]
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed. Please try again.", http.StatusBadRequest)
			return
		}

		client := conf.Client(oauth2.NoContext, tok)
		userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed. Please try again.", http.StatusBadRequest)
			return
		}
		defer userinfo.Body.Close()

		data, _ := ioutil.ReadAll(userinfo.Body)
		u := models.GoogleUser{}
		if err = json.Unmarshal(data, &u); err != nil {
			log.Println(err)
			http.Error(w, "Error marshalling response. Please try again.", http.StatusBadRequest)
			return
		}

		session.Values["user_id"] = u.Sub
		session.Save(r, w)

		// Save or update user here

		http.Redirect(w, r, "/tasks", http.StatusFound)
	}
}

// WebLogin is the login page for the web interface
func WebLogin(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store := sessions.NewCookieStore(env.SessionKey)

		state := RandToken(32)
		session, _ := store.Get(r, "session-name")
		session.Values["state"] = state
		log.Printf("Stored session: %v\n", state)
		session.Save(r, w)

		conf := createOauthConf(env)

		link := conf.AuthCodeURL(state)

		t, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, &viewBag{Link: link})
	}
}

// Index is the homepage of the web application
func Index(_ *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome!")
	}
}

func Web(_ *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("./static/web.html")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(content)
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
