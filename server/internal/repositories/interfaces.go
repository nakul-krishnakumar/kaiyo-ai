package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/models"
)

// UserRepository defines user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Exists(ctx context.Context, email string) (bool, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	UpdateLastLoginAsync(userID uuid.UUID, email string)
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByToken(ctx context.Context, refreshToken string) (*models.Session, error)
	RevokeSessionByToken(ctx context.Context, refreshToken string) error
	RevokeSessionByUserID(ctx context.Context, userID uuid.UUID) error
}

// Repositories aggregates all repositories
type Repositories struct {
	User    UserRepository
	Session SessionRepository
}
