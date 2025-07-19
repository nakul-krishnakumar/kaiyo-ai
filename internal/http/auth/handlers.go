package auth

import (
	"encoding/json"
	"net/http"
)

func NewHandler(ctrl *Controller) *AuthHandler {
	return &AuthHandler{ Controller: ctrl }
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
  	var req struct{ Email, Password string }
  	json.NewDecoder(r.Body).Decode(&req)
	//user := db.VerifyUser(req.Email, req.Password) // your user lookup
	user := struct{ ID, email string }{
		ID: "1",
		email: "dummy@gmail.com",
	}

	// if user == nil {
	// 	http.Error(w, "invalid credentials", http.StatusUnauthorized)
	// 	return
	// }

	access, err := h.Controller.GenerateAccessToken(user.ID, user.email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	refresh, err := h.Controller.GenerateRefreshToken(user.ID, user.email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	// db.StoreRefreshToken(auth.HashToken(refresh), user.ID, time.Now().Add(refreshTTL))
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: refresh,
	})

	json.NewEncoder(w).Encode(map[string]string{"access_token": access})
}