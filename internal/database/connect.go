package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(cfg *Config) (*Database, error) {
	// Building connection config from connection string
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse connection string: " + err.Error())
	}

	// Setting configs for the database pool
	config.MinConns = cfg.Pool.MinConns
	config.MaxConns = cfg.Pool.MaxConns
	config.MaxConnIdleTime = cfg.Pool.MaxConnIdle
	config.MaxConnLifetime = cfg.Pool.MaxConnLife

	// Establishing connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Could not establish connection pool: " + err.Error())
	}

	return &Database{
		Pool: pool,
	}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	if db.Pool != nil {
		slog.Info("Closing database connection pool")
		db.Pool.Close()
	}
}

// Health checks the database connection
func (db *Database) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Stats returns connection pool statistics
func (db *Database) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// WithTx executes a function within a database transaction
func (db *Database) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				slog.Error("Failed to rollback transaction", slog.String("error", rbErr.Error()))
			}
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

//TODO:
/*
	- Setup db repositories
	- db health check api endpoint
	- db health check cron
*/
