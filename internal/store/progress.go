package store

import (
	"context"
	"fmt"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

// DeliverMissionTxParams contains the input parameters for DeliverMissionTx
type DeliverMissionTxParams struct {
	UserID               int64
	MissionID            int64
	NewExperiencePoints  int32
	NewLevel             int32
	WatchPositionSeconds int
}

// DeliverMissionTxResult contains the result of DeliverMissionTx
type DeliverMissionTxResult struct {
	User     *model.User
	Progress *model.UserMissionProgress
}

// DeliverMissionTx delivers a mission reward atomically within a transaction.
// It performs the following operations:
// 1. Gets and locks the progress record to prevent concurrent deliveries
// 2. Re-validates that the mission has not already been delivered (TOCTOU protection)
// 3. Updates user experience and level
// 4. Updates progress status to DELIVERED
//
// This ensures that rewards can only be claimed once, preventing double-reward
// exploits even with concurrent delivery attempts.
func (s *Store) DeliverMissionTx(ctx context.Context, arg DeliverMissionTxParams) (DeliverMissionTxResult, error) {
	var result DeliverMissionTxResult

	err := s.execTx(ctx, func(q *db.Queries) error {
		// Create transaction-aware repositories
		progressRepository := repository.NewProgressRepository(q)
		userRepository := repository.NewUserRepository(q)

		// 1. Get and lock progress record
		// This prevents concurrent delivery attempts
		progress, err := progressRepository.GetByUserAndMission(ctx, arg.UserID, arg.MissionID)
		if err != nil {
			return fmt.Errorf("failed to get progress: %w", err)
		}

		// 2. Re-validate status is not DELIVERED
		// This provides TOCTOU protection - even if pre-transaction check passed,
		// we verify again inside the transaction
		if progress.Status == "DELIVERED" {
			return ErrMissionAlreadyDelivered
		}

		// 3. Update user experience and level
		updatedUser, err := userRepository.UpdateExperience(ctx, arg.UserID, arg.NewExperiencePoints, arg.NewLevel)
		if err != nil {
			return fmt.Errorf("failed to update user experience: %w", err)
		}
		result.User = updatedUser

		// 4. Update progress status to DELIVERED
		updatedProgress, err := progressRepository.Upsert(ctx, arg.UserID, arg.MissionID, "DELIVERED", arg.WatchPositionSeconds)
		if err != nil {
			return fmt.Errorf("failed to update progress: %w", err)
		}
		result.Progress = updatedProgress

		return nil
	})
	if err != nil {
		return DeliverMissionTxResult{}, fmt.Errorf("deliver mission transaction failed: %w", err)
	}

	return result, nil
}
