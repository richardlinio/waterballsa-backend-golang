package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// ErrRefreshTokenNotFound is returned when a refresh token is not found or is invalid
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// RefreshTokenRepository implements refresh token operations using sqlc generated queries.
type RefreshTokenRepository struct {
	queries db.Querier
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
func NewRefreshTokenRepository(queries db.Querier) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		queries: queries,
	}
}

// Create adds a new refresh token to the database
func (r *RefreshTokenRepository) Create(ctx context.Context, jti string, userID int64, expiresAt time.Time) error {
	return r.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		TokenJti: jti,
		UserID:   userID,
		ExpiresAt: pgtype.Timestamp{
			Time:  expiresAt,
			Valid: true,
		},
	})
}

// GetByJTI retrieves a valid (non-revoked, non-expired) refresh token by its JTI
func (r *RefreshTokenRepository) GetByJTI(ctx context.Context, jti string) (*model.RefreshToken, error) {
	dbToken, err := r.queries.GetRefreshToken(ctx, jti)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	token := &model.RefreshToken{
		ID:        dbToken.ID,
		TokenJTI:  dbToken.TokenJti,
		UserID:    dbToken.UserID,
		ExpiresAt: dbToken.ExpiresAt.Time,
		CreatedAt: dbToken.CreatedAt.Time,
	}

	if dbToken.RevokedAt.Valid {
		token.RevokedAt = &dbToken.RevokedAt.Time
	}

	return token, nil
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(ctx context.Context, jti string) error {
	return r.queries.RevokeRefreshToken(ctx, jti)
}

// RevokeAllForUser revokes all refresh tokens for a specific user.
// This method is reserved for future "logout all devices" feature and is not currently exposed
// through any consumer-defined interface.
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	return r.queries.RevokeAllUserRefreshTokens(ctx, userID)
}

// DeleteExpired removes expired refresh tokens from the database
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.queries.DeleteExpiredRefreshTokens(ctx)
}
