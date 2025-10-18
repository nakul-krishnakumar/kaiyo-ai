package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/models"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

// Create a new Session with the details
func (r *SessionRepo) Create(ctx context.Context, session *models.Session) error {
	query := `
	INSERT INTO sessions 
	(user_id, token_hash, expires_at, created_at, last_used_at, device_info)
    VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	session.CreatedAt = now
	session.LastUsedAt = &now

	_, err := r.pool.Exec(ctx, query, session.UserID, session.TokenHash, session.ExpiresAt, session.CreatedAt, session.LastUsedAt, session.DeviceInfo)

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) GetByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	return nil, nil
}
func (r *SessionRepo) RevokeSessionByToken(ctx context.Context, refreshToken string) error {
	return nil
}

func (r *SessionRepo) RevokeSessionByUserID(ctx context.Context, userID uuid.UUID) error {
	return nil
}
