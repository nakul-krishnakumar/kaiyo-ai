package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/nakul-krishnakumar/kaiyo-ai/internal/http/auth"
)

type contextKey string

func Auth(config *auth.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if token == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(token, config.GetAccessSecret())
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), contextKey("userID"), claims.UserID)
			ctx = context.WithValue(ctx, contextKey("email"), claims.Email)
			ctx = context.WithValue(ctx, contextKey("subject"), claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}