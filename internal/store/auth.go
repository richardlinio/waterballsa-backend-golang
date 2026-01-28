package store

import (
	"context"
	"fmt"
	"time"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

// RefreshTokenTxParams contains the input parameters for RefreshTokenTx
type RefreshTokenTxParams struct {
	OldJTI       string
	UserID       int64
	NewJTI       string
	NewExpiresAt time.Time
}

// RefreshTokenTxResult contains the result of RefreshTokenTx
type RefreshTokenTxResult struct {
	Success bool
}

// RefreshTokenTx performs token rotation atomically within a transaction.
// It performs the following operations:
// 1. Re-verifies old token exists and matches userID (protects against concurrent logout)
// 2. Revokes old refresh token
// 3. Creates new refresh token
//
// This ensures that token rotation is atomic - either both revocation and creation
// succeed, or neither happens. This prevents users from being locked out due to
// partial token rotation failures.
func (s *Store) RefreshTokenTx(ctx context.Context, arg RefreshTokenTxParams) (RefreshTokenTxResult, error) {
	var result RefreshTokenTxResult

	err := s.execTx(ctx, func(q *db.Queries) error {
		// Create transaction-aware repository
		refreshTokenRepository := repository.NewRefreshTokenRepository(q)

		// 1. Re-verify old token exists and matches userID
		// This protects against concurrent logout attempts
		oldToken, err := refreshTokenRepository.GetByJTI(ctx, arg.OldJTI)
		if err != nil {
			// Return the repository error directly so service can identify it
			return err
		}

		// Verify user ID matches
		if oldToken.UserID != arg.UserID {
			return ErrTokenUserMismatch
		}

		// 2. Revoke old refresh token
		if err := refreshTokenRepository.Revoke(ctx, arg.OldJTI); err != nil {
			return fmt.Errorf("failed to revoke old token: %w", err)
		}

		// 3. Create new refresh token
		if err := refreshTokenRepository.Create(ctx, arg.NewJTI, arg.UserID, arg.NewExpiresAt); err != nil {
			return fmt.Errorf("failed to create new token: %w", err)
		}

		result.Success = true
		return nil
	})
	if err != nil {
		return RefreshTokenTxResult{}, fmt.Errorf("refresh token transaction failed: %w", err)
	}

	return result, nil
}
