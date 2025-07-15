package chat

import (
	"encoding/json"
	"net/http"
)

type ChatHandler struct {
	Controller *Controller
}

func NewHandler(ctrl *Controller) *ChatHandler {
	return &ChatHandler{ Controller: ctrl }
}

// POST api/v1/chats/
func (h *ChatHandler) PostChat(w http.ResponseWriter, r *http.Request) {
	var reqBody struct { ChatID, UserID, Content string }
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.Content == "" {
		http.Error(w, "content is missing", http.StatusBadRequest)
		return
	} 

	reply, err := h.Controller.SendMessage(reqBody.ChatID, reqBody.UserID, reqBody.Content,)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"reply" : reply,
	}) 
	
	if encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusBadGateway)
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