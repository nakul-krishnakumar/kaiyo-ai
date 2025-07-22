package auth

import (
	"net/http"
)

func New(config *Config) *http.ServeMux {
	ctrl := NewController(config)
	h := NewHandler(ctrl)
	mux := http.NewServeMux()

	dummy := func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement dummy handler logic
	}

	mux.HandleFunc("POST /login", h.EmailLoginHandler)
	mux.HandleFunc("GET /logout", h.EmailLogoutHandler)
	mux.HandleFunc("GET /refresh", h.EmailRefreshHandler)

	mux.HandleFunc("GET /google/callback", dummy)

	return mux
}
