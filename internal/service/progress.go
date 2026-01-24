package service

import (
	"context"
	"errors"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"github.com/linporu/waterballsa-backend-golang/internal/util"
)

const (
	StatusUncompleted = "UNCOMPLETED"
	StatusCompleted   = "COMPLETED"
	StatusDelivered   = "DELIVERED"
)

type progressRepository interface {
	GetByUserAndMission(ctx context.Context, userID, missionID int64) (*model.UserMissionProgress, error)
	Upsert(ctx context.Context, userID, missionID int64, status string, watchPositionSeconds int) (*model.UserMissionProgress, error)
}

type progressMissionRepository interface {
	ListResourcesByMissionID(ctx context.Context, missionID int64) ([]*model.MissionResource, error)
	GetRewardByMissionID(ctx context.Context, missionID int64) (*model.Reward, error)
	GetByID(ctx context.Context, missionID int64) (*model.Mission, int64, error)
}

type progressUserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	UpdateExperience(ctx context.Context, userID int64, experiencePoints, level int32) (*model.User, error)
}

type ProgressService struct {
	progressRepository progressRepository
	missionRepository  progressMissionRepository
	userRepository     progressUserRepository
}

func NewProgressService(
	progressRepository *repository.ProgressRepository,
	missionRepository *repository.MissionRepository,
	userRepository *repository.UserRepository,
) *ProgressService {
	return &ProgressService{
		progressRepository: progressRepository,
		missionRepository:  missionRepository,
		userRepository:     userRepository,
	}
}

// GetProgress retrieves user's progress for a mission
// Returns default progress (UNCOMPLETED, 0 seconds) if mission exists but no progress record
func (s *ProgressService) GetProgress(ctx context.Context, userID, missionID int64) (*model.UserMissionProgress, error) {
	// Get progress (includes mission validation via LEFT JOIN)
	progress, err := s.progressRepository.GetByUserAndMission(ctx, userID, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrProgressNotFound) {
			// Mission doesn't exist
			return nil, apperror.MissionNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	return progress, nil
}

// UpdateProgress updates user's progress for a mission
// Automatically sets status to COMPLETED when watchPositionSeconds >= video duration
func (s *ProgressService) UpdateProgress(ctx context.Context, userID, missionID int64, watchPositionSeconds int) (*model.UserMissionProgress, error) {
	// Validate mission exists and get existing progress
	existingProgress, err := s.progressRepository.GetByUserAndMission(ctx, userID, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrProgressNotFound) {
			return nil, apperror.MissionNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Get video duration from mission resources
	resources, err := s.missionRepository.ListResourcesByMissionID(ctx, missionID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// Calculate total video duration
	var totalDuration int
	for _, resource := range resources {
		if resource.ResourceType == "VIDEO" && resource.DurationSeconds != nil {
			totalDuration += *resource.DurationSeconds
		}
	}

	// Validate watch position
	if watchPositionSeconds < 0 || watchPositionSeconds > totalDuration {
		return nil, apperror.InvalidWatchPosition()
	}

	// Determine status based on watch position
	status := StatusUncompleted
	if watchPositionSeconds >= totalDuration && totalDuration > 0 {
		status = StatusCompleted
	}

	// Preserve COMPLETED status when rewatching (e.g., user watches from middle)
	// Once marked as COMPLETED, it should remain COMPLETED unless explicitly reset
	if existingProgress.Status == StatusCompleted {
		status = StatusCompleted
	}

	// Preserve DELIVERED status even when rewatching
	if existingProgress.Status == StatusDelivered {
		status = StatusDelivered
	}

	// Upsert progress
	progress, err := s.progressRepository.Upsert(ctx, userID, missionID, status, watchPositionSeconds)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return progress, nil
}

// DeliverMission delivers a completed mission and grants experience points
// Business logic:
// 1. Verify mission exists and get mission type
// 2. Get user progress (must exist)
// 3. Check progress status:
//   - VIDEO missions: must be COMPLETED
//   - Other missions: can be UNCOMPLETED or COMPLETED
//   - Any mission: cannot be DELIVERED (409 Conflict)
//
// 4. Get mission reward (EXPERIENCE type)
// 5. Update user experience and level
// 6. Update progress status to DELIVERED
func (s *ProgressService) DeliverMission(ctx context.Context, userID, missionID int64) (*model.MissionDeliveryResult, error) {
	// Get mission details to check mission type
	mission, _, err := s.missionRepository.GetByID(ctx, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrMissionNotFound) {
			return nil, apperror.MissionNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Get progress record
	progress, err := s.progressRepository.GetByUserAndMission(ctx, userID, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrProgressNotFound) {
			return nil, apperror.MissionNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Check if already delivered
	if progress.Status == StatusDelivered {
		return nil, apperror.MissionAlreadyDelivered()
	}

	// Validate status based on mission type
	if mission.Type == "VIDEO" && progress.Status != StatusCompleted {
		return nil, apperror.MissionNotCompleted()
	}
	// Non-video missions can be delivered in UNCOMPLETED or COMPLETED status

	// Get mission reward
	reward, err := s.missionRepository.GetRewardByMissionID(ctx, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrRewardNotFound) {
			// No reward configured, set to 0
			reward = &model.Reward{RewardValue: 0}
		} else {
			return nil, apperror.DatabaseError(err)
		}
	}

	// Get current user data
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperror.Unauthorized()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Calculate new experience and level
	//nolint:gosec // reward.RewardValue is validated to be within reasonable bounds by database schema
	newExp := user.ExperiencePoints + int32(reward.RewardValue)
	newLevel := util.CalculateLevelFromExperience(newExp)

	// Update user experience and level
	_, err = s.userRepository.UpdateExperience(ctx, userID, newExp, newLevel)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// Update progress status to DELIVERED
	_, err = s.progressRepository.Upsert(ctx, userID, missionID, StatusDelivered, progress.WatchPositionSeconds)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &model.MissionDeliveryResult{
		Message:          "任務交付成功",
		ExperienceGained: reward.RewardValue,
		TotalExperience:  newExp,
		CurrentLevel:     newLevel,
	}, nil
}
