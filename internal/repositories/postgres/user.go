package postgres

import (
	"context"
	"fmt"
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

// Create implements repositories.UserRepository
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, password, name, google_id, twitter_id, 
                          email_verified, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.Name,
		user.GoogleID, user.TwitterID, user.EmailVerified, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID implements repositories.UserRepository
func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT id, email, password, name, google_id, twitter_id,
               email_verified, is_active, created_at, updated_at, last_login_at
        FROM users WHERE id = $1 AND is_active = true`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
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

// GetByEmail implements repositories.UserRepository
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password, name, google_id, twitter_id,
               email_verified, is_active, created_at, updated_at, last_login_at
        FROM users WHERE email = $1 AND is_active = true`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
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

// GetByGoogleID implements repositories.UserRepository
func (r *UserRepo) GetByGoogleID(ctx context.Context, googleId string) (*models.User, error) {
	query := `
        SELECT id, email, password, name, google_id, twitter_id,
               email_verified, is_active, created_at, updated_at, last_login_at
        FROM users WHERE google_id = $1 AND is_active = true`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, googleId).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
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

// GetByTwitterID implements repositories.UserRepository
func (r *UserRepo) GetByTwitterID(ctx context.Context, twitterId string) (*models.User, error) {
	query := `
        SELECT id, email, password, name, google_id, twitter_id,
               email_verified, is_active, created_at, updated_at, last_login_at
        FROM users WHERE twitter_id = $1 AND is_active = true`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, twitterId).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
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

// Update implements repositories.UserRepository
func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users 
        SET email = $2, name = $3, google_id = $4, twitter_id = $5,
            email_verified = $6, updated_at = $7
        WHERE id = $1 AND is_active = true`

	user.UpdatedAt = time.Now()

	result, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.GoogleID, user.TwitterID,
		user.EmailVerified, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete implements repositories.UserRepository
func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete by setting is_active to false
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1 AND is_active = true`

	result, err := r.pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepo) VerifyUser(ctx context.Context, email, password string) (*models.User, error) {
	query := `
	SELECT id, email, password, name, google_id, twitter_id,
			email_verified, is_active, created_at, updated_at, last_login_at
	FROM users WHERE email = $1 AND password = $2`

	user := &models.User{}
		err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.GoogleID, &user.TwitterID, &user.EmailVerified, &user.IsActive,
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