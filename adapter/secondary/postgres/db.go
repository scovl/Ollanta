package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgxpool.Pool and provides lightweight helpers.
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a DB by opening a pgxpool connection to the given DSN.
func New(dsn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	return &DB{Pool: pool}, nil
}

// Health pings the database and returns an error if unreachable.
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
