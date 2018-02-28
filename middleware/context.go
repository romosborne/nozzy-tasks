package middleware

import (
	"context"
	"net/http"

	"github.com/romosborne/nozzy-tasks/services"
)

const (
	SQLServiceKey    = "SQLServiceKey"
	OauthClientIDKey = "OauthClientIDKey"
)

func AddContext(next http.Handler, sql *services.SQL, oauthClientID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), SQLServiceKey, sql)
		ctx = context.WithValue(ctx, OauthClientIDKey, oauthClientID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
