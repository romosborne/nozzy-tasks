package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nozzy-tasks/models"

	"github.com/futurenda/google-auth-id-token-verifier"
	"github.com/gorilla/sessions"
)

func WebValidate(env *models.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := sessions.NewCookieStore(env.SessionKey)
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Please login", http.StatusInternalServerError)
			return
		}

		x := session.Values["user_id"]

		if x == nil {
			http.Error(w, "Please login", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ApiValidate(env *models.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authorizationHeader := req.Header.Get("authorization")
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

			ctx := context.WithValue(req.Context(), fmt.Sprintf("%s_id", env.ContextKey), claimSet.Sub)
			ctx = context.WithValue(ctx, fmt.Sprintf("%s_email", env.ContextKey), claimSet.Email)

			next.ServeHTTP(w, req.WithContext(ctx))
		}
	})
}
