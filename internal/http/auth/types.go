package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Controller struct {
	auth *Config
}

type JWTCustomClaims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type UserReq struct {
	UserName string
	Email    string
	Password string
}

type Handler struct {
	Controller *Controller
}
