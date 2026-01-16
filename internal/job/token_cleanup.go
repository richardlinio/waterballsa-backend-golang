package job

import (
	"context"
	"fmt"
	"time"

	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/scheduler"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
)

// expiredAccessTokenCleaner is used by TokenCleanup job to clean expired access tokens
type expiredAccessTokenCleaner interface {
	DeleteExpired(ctx context.Context) error
}

// expiredRefreshTokenCleaner is used by TokenCleanup job to clean expired refresh tokens
type expiredRefreshTokenCleaner interface {
	DeleteExpired(ctx context.Context) error
}

// TokenCleanup encapsulates the token cleanup job
type TokenCleanup struct {
	accessTokenCleaner  expiredAccessTokenCleaner
	refreshTokenCleaner expiredRefreshTokenCleaner
	interval            time.Duration
}

// NewTokenCleanup creates a new TokenCleanup job
func NewTokenCleanup(
	accessTokenRepository *repository.AccessTokenRepository,
	refreshTokenRepository *repository.RefreshTokenRepository,
	interval time.Duration,
) *TokenCleanup {
	return &TokenCleanup{
		accessTokenCleaner:  accessTokenRepository,
		refreshTokenCleaner: refreshTokenRepository,
		interval:            interval,
	}
}

// ToSchedulerJob converts TokenCleanup to a scheduler.Job
func (j *TokenCleanup) ToSchedulerJob() scheduler.Job {
	return scheduler.Job{
		Name:     "token-cleanup",
		Interval: j.interval,
		Run: func(ctx context.Context) error {
			if err := j.accessTokenCleaner.DeleteExpired(ctx); err != nil {
				return fmt.Errorf("failed to cleanup expired access tokens: %w", err)
			}
			if err := j.refreshTokenCleaner.DeleteExpired(ctx); err != nil {
				return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
			}
			return nil
		},
	}
}
