package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nozzy-tasks/models"

	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
)

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(models.Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next.ServeHTTP(w, req)
				} else {
					json.NewEncoder(w).Encode(models.Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(models.Exception{Message: "An authorization header is required"})
		}
	})
}
