package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/models"
)

// UserRepository defines user data operations
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetByGoogleID(ctx context.Context, googleId string) (*models.User, error)
    GetByTwitterID(ctx context.Context, twitterId string) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Repositories aggregates all repositories
type Repositories struct {
    User UserRepository
}
