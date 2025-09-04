package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func NewController(config *Config) *Controller {
	return &Controller{
		config,
	}
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading authentication secrets" + err.Error())
	}

	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	if len(accessSecret) == 0 {
		slog.Error("could not load access secret")
		os.Exit(1)
	}

	refreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	if len(refreshSecret) == 0 {
		slog.Error("could not load refresh secret")
		os.Exit(1)
	}

	accessTTL, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		slog.Error("could not load access TTL " + err.Error())
	}

	refreshTTL, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		slog.Error("could not load refresh TTL " + err.Error())
	}

	return &Config{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
	}
}

// AccessSecret Getter
func (c *Config) GetAccessSecret() []byte {
	return c.AccessSecret
}

func NewJWTCustomClaims(userID uuid.UUID, email string, TTL time.Duration) *JWTCustomClaims {
	return &JWTCustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func (c *Controller) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := NewJWTCustomClaims(userID, email, c.auth.AccessTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(c.auth.AccessSecret)
	if err != nil {
		return "", fmt.Errorf("error generating access token %s", err)
	}

	return ss, nil
}

func (c *Controller) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	claims := NewJWTCustomClaims(userID, email, c.auth.RefreshTTL)
	claims.ID = uuid.NewString()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(c.auth.RefreshSecret)
	if err != nil {
		return "", fmt.Errorf("error generating refresh token %s", err)
	}

	return ss, nil
}

func ValidateToken(tokenStr string, secret []byte) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr,
		&JWTCustomClaims{}, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return secret, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
