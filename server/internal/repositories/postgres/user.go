package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create a new User with details
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users 
		(id, email, password, first_name, last_name, username, google_id, twitter_id, email_verified, is_active, created_at, updated_at, last_login_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastLoginAt = now

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.UserName, user.GoogleID, user.TwitterID, user.EmailVerified, user.IsActive,
		user.CreatedAt, user.UpdatedAt, user.LastLoginAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Check if user with the passed email exists
func (r *UserRepo) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// GetByEmail implements repositories.UserRepository
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password, first_name, last_name, google_id, twitter_id,
               email_verified, is_active, created_at, updated_at, last_login_at
        FROM users WHERE email = $1 AND is_active = true`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
	UPDATE TABLE user 
	SET last_login_at = NOW()
	WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

func (r *UserRepo) UpdateLastLoginAsync(userID uuid.UUID, email string) {
	go func() {
		// Panic recovery for goroutine safety
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic in UpdateLastLoginAsync",
					slog.Any("panic", r),
					slog.String("user_id", userID.String()))
			}
		}()

		// Create background context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Retry logic with exponential backoff
		maxRetries := 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := r.UpdateLastLogin(ctx, userID); err != nil {
				slog.Error("Failed to update last login",
					slog.String("error", err.Error()),
					slog.String("user_id", userID.String()),
					slog.String("email", email),
					slog.Int("attempt", attempt))

				if attempt < maxRetries {
					// Exponential backoff: 1s, 4s, 9s
					backoff := time.Duration(attempt*attempt) * time.Second

					select {
					case <-time.After(backoff):
						continue
					case <-ctx.Done():
						slog.Error("Context cancelled during retry backoff",
							slog.String("user_id", userID.String()))
						return
					}
				}

				// All retries failed
				slog.Error("All retries failed for last login update",
					slog.String("user_id", userID.String()),
					slog.String("email", email))
				return
			}

			// Success - log and exit
			slog.Debug("Last login updated successfully",
				slog.String("user_id", userID.String()),
				slog.Int("attempt", attempt))
			return
		}
	}()
}
