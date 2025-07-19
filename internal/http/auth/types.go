package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Controller struct {
	auth *AuthConfig
}

type JWTCustomClaims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type UserReq struct {
	UserName string
	Email string
	Password string
}

type AuthHandler struct {
	Controller *Controller
}
