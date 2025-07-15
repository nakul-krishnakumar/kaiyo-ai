package chat

import (
	"net/http"
)

func NewMux() *http.ServeMux {
	ctrl := NewController()
	h := NewHandler(ctrl)

	m := http.NewServeMux()
	m.HandleFunc("POST /{$}", h.PostChat) 				    // api/v1/chats/
	m.HandleFunc("GET /history/{chatID}", h.GetHistory)     // api/v1/chats/history/{chatID}

	return m
}