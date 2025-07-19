package auth

import "net/http"


func New(auth *AuthConfig) *http.ServeMux {
	ctrl := NewController(auth)
	h := NewHandler(ctrl)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", h.LoginHandler)

	return mux
}