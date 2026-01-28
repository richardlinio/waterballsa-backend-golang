package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
)

// Store provides all functions to run database queries and transactions
type Store struct {
	*db.Queries
	pool *pgxpool.Pool
}

// New creates a new Store instance
func New(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: db.New(pool),
		pool:    pool,
	}
}

// execTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic and re-panic
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// Rollback on error
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
			}
		} else {
			// Commit on success
			err = tx.Commit(ctx)
		}
	}()

	q := db.New(tx)
	err = fn(q)
	return err
}
