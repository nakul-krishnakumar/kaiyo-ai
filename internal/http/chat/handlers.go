package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatHandler struct {
	Controller *Controller
}

type UserInput struct {
	ChatID string
	UserID string
	Content string
}

func NewHandler(ctrl *Controller) *ChatHandler {
	return &ChatHandler{ Controller: ctrl }
}

// POST api/v1/chats/
func (h *ChatHandler) PostChat(w http.ResponseWriter, r *http.Request) {

	chunkCh := make(chan string)
	done := make(chan error, 1)

	ctx := r.Context()

	// Immediate streaming without waiting for the entire response
	flusher, ok := w.(http.Flusher); 
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	userInput :=  UserInput{}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userInput.Content == "" {
		http.Error(w, "content is missing", http.StatusBadRequest)
		return
	} 

	go func() {
		done <- h.Controller.StreamMessage(ctx, userInput, chunkCh)	
	}()

	// listen to both done and chunkCh channels at the same time
	for { // inifinite loop
		select { // listen from the channel which gives output first
		case <-ctx.Done():
			return
	
		case err := <-done:
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case chunk, ok := <-chunkCh:
			if !ok { return }
			fmt.Fprintf(w, "data: %s\n\n", chunk) // buffer
			flusher.Flush() // returns the buffer

			fmt.Print(chunk)
		}
  	}
}

func (h *ChatHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chatID")
	
	if chatID == "" {
		http.Error(w, "parameter chatID is missing", http.StatusBadRequest)
		return
	}

	msgs, err := h.Controller.GetHistory(chatID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(msgs)
}