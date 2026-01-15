package job

import (
	"context"
	"fmt"
	"time"

	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/scheduler"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
)

// TokenCleanup encapsulates the token cleanup job
type TokenCleanup struct {
	accessTokenRepository  repository.AccessTokenRepository
	refreshTokenRepository repository.RefreshTokenRepository
	interval               time.Duration
}

// NewTokenCleanup creates a new TokenCleanup job
func NewTokenCleanup(
	accessTokenRepository repository.AccessTokenRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	interval time.Duration,
) *TokenCleanup {
	return &TokenCleanup{
		accessTokenRepository:  accessTokenRepository,
		refreshTokenRepository: refreshTokenRepository,
		interval:               interval,
	}
}

// ToSchedulerJob converts TokenCleanup to a scheduler.Job
func (j *TokenCleanup) ToSchedulerJob() scheduler.Job {
	return scheduler.Job{
		Name:     "token-cleanup",
		Interval: j.interval,
		Run: func(ctx context.Context) error {
			if err := j.accessTokenRepository.DeleteExpired(ctx); err != nil {
				return fmt.Errorf("failed to cleanup expired access tokens: %w", err)
			}
			if err := j.refreshTokenRepository.DeleteExpired(ctx); err != nil {
				return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
			}
			return nil
		},
	}
}
