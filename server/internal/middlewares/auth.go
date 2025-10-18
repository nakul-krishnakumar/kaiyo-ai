package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nakul-krishnakumar/kaiyo-ai/internal/http/auth"
)

type contextKey string

func Auth(config *auth.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing authorization header",
				})
				return
			}

			claims, err := auth.ValidateToken(token, config.GetAccessSecret())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "invalid token",
					"error":   err.Error(),
				})
				return
			}

			ctx := context.WithValue(r.Context(), contextKey("userID"), claims.UserID)
			ctx = context.WithValue(ctx, contextKey("subject"), claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
