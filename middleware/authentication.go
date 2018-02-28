package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/romosborne/nozzy-tasks/models"
	"github.com/romosborne/nozzy-tasks/services"
)

const UserContextKey = "userContext"

// APIValidate validates incoming calls to the api and populates the user object in the context
func APIValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader == "" {
			json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			sqlService := getSQLService(req)

			user, err := sqlService.GetUserFromAuthToken(authToken)

			if err != nil {
				json.NewEncoder(w).Encode(models.Exception{Message: "Invalid auth token"})
				return
			}

			ctx := context.WithValue(req.Context(), UserContextKey, user)

			next.ServeHTTP(w, req.WithContext(ctx))
		}
	})
}

func getSQLService(r *http.Request) *services.SQL {
	return r.Context().Value(SQLServiceKey).(*services.SQL)
}
