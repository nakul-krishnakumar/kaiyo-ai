package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/repositories"
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
	UserID uuid.UUID `json:"userID"`
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
	Repo *repositories.Repositories
	Validator Validator
}

type SignInRequest struct {
	FirstName string `json:"firstName" mapstructure:"first_name"`
	LastName  string `json:"lastName" mapstructure:"last_name"`
	Email     string `json:"email" mapstructure:"email"`
	Password  string `json:"password" mapstructure:"password"`
}

type LogInRequest struct {
	Email    string `json:"email" mapstructure:"email"`
	Password string `json:"password" mapstructure:"password"`
}

type Validator interface {
    ValidateSignInRequest(req SignInRequest) ValidationErrors
    ValidateLoginRequest(req LogInRequest) ValidationErrors
}

// ValidationErrors type for structured error handling
type ValidationErrors map[string]string

type validator struct{}