package chat

import (
	"net/http"

	mw "github.com/nakul-krishnakumar/kaiyo-ai/internal/middlewares"
)

func New() *http.ServeMux {
	ctrl := NewController()
	h := NewHandler(ctrl)

	mux := http.NewServeMux()
	mux.Handle("POST /{$}", mw.SSEHandler(http.HandlerFunc(h.PostChat)))//(sse output) api/v1/chats/
	mux.HandleFunc("GET /history/{chatID}", h.GetHistory)     // api/v1/chats/history/{chatID}

	return mux
}