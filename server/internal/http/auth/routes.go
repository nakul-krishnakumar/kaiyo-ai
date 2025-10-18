package auth

import (
	"net/http"

	"github.com/nakul-krishnakumar/kaiyo-ai/internal/repositories"
)

func New(config *Config, repo *repositories.Repositories) *http.ServeMux {
	ctrl := NewController(config)
	val := NewValidator()
	h := NewHandler(ctrl, repo, val)
	mux := http.NewServeMux()

	dummy := func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement dummy handler logic
	}

	mux.HandleFunc("POST /login", h.EmailLoginHandler)
	mux.HandleFunc("POST /signin", h.EmailSignInHandler)
	mux.HandleFunc("GET /logout", h.EmailLogoutHandler)
	mux.HandleFunc("GET /refresh", h.EmailRefreshHandler)

	mux.HandleFunc("GET /google/callback", dummy)

	return mux
}
