package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/models"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCostFactor = 14
)

func NewHandler(ctrl *Controller, repo *repositories.Repositories, val Validator) *Handler {
	return &Handler{Controller: ctrl, Repo: repo, Validator: val}
}

func (h *Handler) SendError(w http.ResponseWriter, err error, msg any) {
	newErr := json.NewEncoder(w).Encode(map[string]any{
		"error":   err.Error(),
		"message": msg,
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

func (h *Handler) GetDeviceInfo(r *http.Request) *models.DeviceInfo {
	userAgent := r.UserAgent()

	return &models.DeviceInfo{
		UserAgent:  userAgent,
		IPAddress:  h.getClientIP(r),
		Language:   h.getLanguage(r),
		Platform:   h.parsePlatform(userAgent),
		Browser:    h.parseBrowser(userAgent),
		OS:         h.parseOS(userAgent),
		DeviceType: h.parseDeviceType(userAgent),
		AppVersion: r.Header.Get("X-App-Version"), // Custom header from mobile apps
		Country:    "",                            // You can add GeoIP lookup later
		LastSeenAt: time.Now(),
	}
}

// Get real client IP (handles proxies and load balancers)
func (h *Handler) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Format: "client-ip, proxy1-ip, proxy2-ip"
		ips := strings.Split(xff, ",")
		clientIP := strings.TrimSpace(ips[0])
		if clientIP != "" {
			return clientIP
		}
	}

	// Check X-Real-IP header (Nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Check X-Forwarded header
	if xf := r.Header.Get("X-Forwarded"); xf != "" {
		return xf
	}

	// Fall back to RemoteAddr (remove port if present)
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// ✅ Get preferred language from Accept-Language header
func (h *Handler) getLanguage(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		return "en"
	}

	// Parse "en-US,en;q=0.9,es;q=0.8" -> "en-US"
	if idx := strings.Index(lang, ","); idx != -1 {
		lang = lang[:idx]
	}

	// Remove quality values like ";q=0.9"
	if idx := strings.Index(lang, ";"); idx != -1 {
		lang = lang[:idx]
	}

	return strings.TrimSpace(lang)
}

// ✅ Parse platform from User-Agent
func (h *Handler) parsePlatform(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os"):
		return "macOS"
	case strings.Contains(ua, "linux") && !strings.Contains(ua, "android"):
		return "Linux"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod"):
		return "iOS"
	case strings.Contains(ua, "freebsd"):
		return "FreeBSD"
	default:
		return "Unknown"
	}
}

// ✅ Parse browser from User-Agent
func (h *Handler) parseBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/"):
		return "Edge"
	case strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg"):
		return "Chrome"
	case strings.Contains(ua, "firefox/"):
		return "Firefox"
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome"):
		return "Safari"
	case strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/"):
		return "Opera"
	case strings.Contains(ua, "curl/"):
		return "cURL"
	case strings.Contains(ua, "postman"):
		return "Postman"
	default:
		return "Unknown"
	}
}

// ✅ Parse OS version from User-Agent
func (h *Handler) parseOS(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows nt 10.0"):
		return "Windows 10"
	case strings.Contains(ua, "windows nt 6.3"):
		return "Windows 8.1"
	case strings.Contains(ua, "windows nt 6.1"):
		return "Windows 7"
	case strings.Contains(ua, "mac os x"):
		// Extract version like "10_15_7"
		if start := strings.Index(ua, "mac os x "); start != -1 {
			versionStart := start + 9
			if end := strings.Index(ua[versionStart:], ")"); end != -1 {
				version := ua[versionStart : versionStart+end]
				version = strings.ReplaceAll(version, "_", ".")
				return "macOS " + version
			}
		}
		return "macOS"
	case strings.Contains(ua, "android"):
		// Extract Android version
		if start := strings.Index(ua, "android "); start != -1 {
			versionStart := start + 8
			if end := strings.Index(ua[versionStart:], ";"); end != -1 {
				version := ua[versionStart : versionStart+end]
				return "Android " + version
			}
		}
		return "Android"
	case strings.Contains(ua, "iphone os"):
		// Extract iOS version
		if start := strings.Index(ua, "os "); start != -1 {
			versionStart := start + 3
			if end := strings.Index(ua[versionStart:], " "); end != -1 {
				version := ua[versionStart : versionStart+end]
				version = strings.ReplaceAll(version, "_", ".")
				return "iOS " + version
			}
		}
		return "iOS"
	case strings.Contains(ua, "ubuntu"):
		return "Ubuntu"
	case strings.Contains(ua, "debian"):
		return "Debian"
	case strings.Contains(ua, "centos"):
		return "CentOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}

// ✅ Parse device type from User-Agent
func (h *Handler) parseDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "mobile"):
		return "mobile"
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return "tablet"
	case strings.Contains(ua, "tv") || strings.Contains(ua, "smart-tv"):
		return "tv"
	case strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider"):
		return "bot"
	case strings.Contains(ua, "curl") || strings.Contains(ua, "postman"):
		return "api_client"
	default:
		return "desktop"
	}
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
		w.WriteHeader(http.StatusNotAcceptable)
		h.SendError(w, fmt.Errorf("Validation failed"), err)
		return
	}

	exists, err := h.Repo.User.Exists(r.Context(), req.Email)
	if exists {
		w.WriteHeader(http.StatusConflict)
		h.SendError(w, fmt.Errorf("signin request failed"), "user already exists")
		return
	}

	if err != nil {
		slog.Error("could not check database for duplicate entries\n", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCostFactor)

	if err != nil {
		slog.Error("could not hash the password", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Password:  string(passwordHash),

		EmailVerified: false,
		IsActive:      true,
	}

	if err = h.Repo.User.Create(r.Context(), &user); err != nil {
		slog.Error("could not create new user", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	access, err := h.Controller.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate access token", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	refresh, err := h.Controller.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate refresh token", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, refresh)

	// hashedToken := HashToken(refresh)
	// session := &models.Session{
	// 	UserID: user.ID,
	// 	TokenHash: hashedToken,
	// 	ExpiresAt: time.Now().Add(h.Controller.auth.RefreshTTL),
	// 	CreatedAt: time.Now(),
	// 	// DeviceInfo: ,
	// }

	// err := h.Repo.Session.Create(r.Context(), u)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]any{
		"message":      "Account created successfully",
		"access_token": access,
		"expires_in":   int64(h.Controller.auth.AccessTTL / time.Second),
		"user": map[string]any{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"createdAt": user.CreatedAt,
		},
	})

	if err != nil {
		slog.Error("could not encode sign in response", slog.String("error", err.Error()))
	}
}

func (h *Handler) EmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LogInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validator.ValidateLoginRequest(req); err.HasErrors() {
		w.WriteHeader(http.StatusBadRequest)
		h.SendError(w, fmt.Errorf("Validation failed"), err)
		return
	}

	user, err := h.Repo.User.GetByEmail(r.Context(), req.Email)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
			"name":       user.FirstName + " " + user.LastName,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})

	if err != nil {
		slog.Error("Could not return access token", slog.String("error", err.Error()))
		return
	}

	h.Repo.User.UpdateLastLoginAsync(user.ID, user.Email)
	slog.Info("Login successful", slog.String("email", user.Email))
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
