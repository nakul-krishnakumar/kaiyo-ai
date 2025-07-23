package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nakul-krishnakumar/kaiyo-ai/internal/repositories/postgres"
)

// NewRepositories creates all repository implementations
func NewRepositories(pool *pgxpool.Pool) *Repositories {
    return &Repositories{
        User: postgres.NewUserRepo(pool),
		// Add other repos here
    }
}