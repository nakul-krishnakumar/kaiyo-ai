package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/repositories"
)
 
func NewHandler(ctrl *Controller, repo *repositories.Repositories, val Validator) *Handler {
	return &Handler{Controller: ctrl, Repo: repo, Validator: val}
}

func (h *Handler) SendError(w http.ResponseWriter, err error, msg any) {
	newErr := json.NewEncoder(w).Encode(map[string]any{
		"error" : err.Error(),
		"message" : msg,
	})

	if newErr != nil {
		slog.Error("Could not encode response", slog.String("error", newErr.Error()))
	}
}

func (h *Handler) isProd() bool {
	if err := godotenv.Load(); err != nil {
		return true
	}

	return os.Getenv("ENVIRONMENT") == "production"
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, refresh string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Secure:   h.isProd(),
		SameSite: http.SameSiteStrictMode, // CSRF protection
		Expires:  time.Now().Add(h.Controller.auth.RefreshTTL),
		MaxAge:   int(h.Controller.auth.RefreshTTL / time.Second),
		// Secure: true, // use it in production, secure: true lets only https reqs to use the cookie
	})
}


func (h *Handler) EmailSignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validator.ValidateSignInRequest(req); err.HasErrors() {
		w.WriteHeader(http.StatusBadRequest)
		h.SendError(w, fmt.Errorf("Validation failed"), err)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message" : "successfully signed in",
	}); err != nil {
		slog.Error("could not encode sign in response", slog.String("error", err.Error()))
	}
}

func (h *Handler) EmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LogInRequest	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	user, err := h.Repo.User.GetByEmail(r.Context(), req.Email)
	fmt.Println("USER ", user) //TODO remove !!!!
	if err != nil {
		http.Error(w, "Invalid credentials: " + err.Error(), http.StatusUnauthorized)
		return
	}

	access, err := h.Controller.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refresh, err := h.Controller.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// db.StoreRefreshToken(auth.HashToken(refresh), user.ID, time.Now().Add(refreshTTL))
	h.setRefreshTokenCookie(w, refresh)

	err = json.NewEncoder(w).Encode(map[string]any{
		"access_token": access,
		"expires_in":   int64(h.Controller.auth.AccessTTL / time.Second),
		"user": map[string]any{
			"name":      user.Name,
			"email":     user.Email,
			"created_at": user.CreatedAt,
		},
	})

	if err != nil {
		slog.Error("Could not return access token", slog.String("error", err.Error()))
	}
}

func (h *Handler) EmailLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if errors.Is(err, http.ErrNoCookie) || cookie == nil {
		slog.Info("Logout attempt with no refresh token")

		// Still return success (idempotent operation)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{
			"message": "Already logged out",
		})

		if err != nil {
			slog.Error("Could not log out user", slog.String("error", err.Error()))
		}

		return
	}

	if err == nil && cookie.Value != "" {
		//TODO : Invalidate token in db
		slog.Info("Refresh token invalidated")
	}

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   h.isProd(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})

	if err != nil {
		slog.Error("Could not respond to logout request", slog.String("error", err.Error()))
	}
}

func (h *Handler) EmailRefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		h.SendError(w, err, "Refresh token does not exists")
		return
	}

	// TODO: Validate refresh token from database
	// claims, err := h.Controller.ValidateRefreshToken(cookie.Value)
	// if err != nil {
	//     http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
	//     return
	// }

	// For now, validate the JWT directly
	claims, err := ValidateToken(cookie.Value, h.Controller.auth.RefreshSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.SendError(w, err, "Invalid token")
		return
	}

	// Generate new access token
	access, err := h.Controller.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		slog.Error("Failed to generate new access token", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Optionally generate new refresh token (token rotation)
	refresh, err := h.Controller.GenerateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		slog.Error("Failed to generate new refresh token", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: Update refresh token in database
	// h.Controller.ReplaceRefreshToken(oldToken, refresh, claims.UserID)

	// Update cookie with new refresh token
	h.setRefreshTokenCookie(w, refresh)

	// Return new access token
	err = json.NewEncoder(w).Encode(map[string]string{
		"message":      "refreshed token successfully",
		"access_token": access,
	})

	if err != nil {
		slog.Error("Could not respond to refresh request", slog.String("error", err.Error()))
		return
	}

	slog.Info("Token refreshed", slog.String("user_id", claims.UserID.String()))
}
