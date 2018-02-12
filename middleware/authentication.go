package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/romosborne/nozzy-tasks/models"
)

// APIValidate validates incoming calls to the api and populates the user object in the context
func APIValidate(env *models.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader == "" {
			json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			user, err := env.Db.GetUserFromAuthToken(authToken)

			if err != nil {
				json.NewEncoder(w).Encode(models.Exception{Message: "Invalid auth token"})
				return
			}

			ctx := context.WithValue(req.Context(), env.ContextKey, user)

			next.ServeHTTP(w, req.WithContext(ctx))
		}
	})
}
