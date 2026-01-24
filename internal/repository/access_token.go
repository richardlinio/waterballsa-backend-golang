package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
)

// AccessTokenRepository implements access token blacklist operations using sqlc generated queries.
type AccessTokenRepository struct {
	queries db.Querier
}

// NewAccessTokenRepository creates a new instance of AccessTokenRepository
func NewAccessTokenRepository(queries db.Querier) *AccessTokenRepository {
	return &AccessTokenRepository{
		queries: queries,
	}
}

// Invalidate adds a token JTI to the blacklist
func (r *AccessTokenRepository) Invalidate(ctx context.Context, jti string, userID int64, expiresAt time.Time) error {
	return r.queries.InvalidateToken(ctx, db.InvalidateTokenParams{
		TokenJti: jti,
		UserID:   userID,
		ExpiresAt: pgtype.Timestamp{
			Time:  expiresAt,
			Valid: true,
		},
	})
}

// IsInvalidated checks if a token JTI is in the blacklist
func (r *AccessTokenRepository) IsInvalidated(ctx context.Context, jti string) (bool, error) {
	return r.queries.IsTokenInvalidated(ctx, jti)
}

// DeleteExpired removes expired tokens from the blacklist
func (r *AccessTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredTokens(ctx)
}
