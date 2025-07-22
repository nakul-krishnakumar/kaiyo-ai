package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func NewHandler(ctrl *Controller) *Handler {
	return &Handler{ Controller: ctrl }
}

func (h *Handler) isProd() bool {
	if err := godotenv.Load(); err != nil {
		return true
	}

	return os.Getenv("ENVIRONMENT") == "production"
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, refresh string) {
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: refresh,
		HttpOnly: true,
		Secure: h.isProd(),
		SameSite: http.SameSiteStrictMode, // CSRF protection
		Expires: time.Now().Add(h.Controller.auth.RefreshTTL),
		MaxAge: int(h.Controller.auth.RefreshTTL / time.Second),
		// Secure: true, // use it in production, secure: true lets only https reqs to use the cookie
	})
}

func (h *Handler) EmailLoginHandler(w http.ResponseWriter, r *http.Request) {

  	var req struct{ Email, Password string }
  	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
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
		return
	}

	refresh, err := h.Controller.GenerateRefreshToken(user.ID, user.email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// db.StoreRefreshToken(auth.HashToken(refresh), user.ID, time.Now().Add(refreshTTL))
	h.setRefreshTokenCookie(w, refresh)

	json.NewEncoder(w).Encode(map[string]any{
		"access_token": access,
		"expires_in": int64(h.Controller.auth.AccessTTL / time.Second),
	})
}

func (h *Handler) EmailLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if errors.Is(err, http.ErrNoCookie) || cookie == nil {
        slog.Info("Logout attempt with no refresh token")
        
        // Still return success (idempotent operation)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Already logged out",
        })

        return
	}

	if err == nil && cookie.Value != "" {
		//TODO : Invalidate token in db
		slog.Info("Refresh token invalidated")
	}

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: "",
		HttpOnly: true,
		Secure: h.isProd(),
		SameSite: http.SameSiteStrictMode,
		MaxAge: -1,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message" : "Logged out successfully",
	})
}

func (h *Handler) EmailRefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error" : err.Error(),
		})
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
		json.NewEncoder(w).Encode(map[string]string{
			"error" : err.Error(),
		})
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
    json.NewEncoder(w).Encode(map[string]string{
		"message" : "refreshed token successfully",
		"access_token" : access,
	})

    slog.Info("Token refreshed", slog.String("user_id", claims.UserID))
}
