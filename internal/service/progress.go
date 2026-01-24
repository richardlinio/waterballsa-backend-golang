package service

import (
	"context"
	"errors"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
)

type progressRepository interface {
	GetByUserAndMission(ctx context.Context, userID, missionID int64) (*model.UserMissionProgress, error)
	Upsert(ctx context.Context, userID, missionID int64, status string, watchPositionSeconds int) (*model.UserMissionProgress, error)
}

type ProgressService struct {
	progressRepository progressRepository
	missionRepository  missionRepository
}

func NewProgressService(
	progressRepository *repository.ProgressRepository,
	missionRepository *repository.MissionRepository,
) *ProgressService {
	return &ProgressService{
		progressRepository: progressRepository,
		missionRepository:  missionRepository,
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
	status := "UNCOMPLETED"
	if watchPositionSeconds >= totalDuration && totalDuration > 0 {
		status = "COMPLETED"
	}

	// Preserve DELIVERED status even when rewatching
	if existingProgress.Status == "DELIVERED" {
		status = "DELIVERED"
	}

	// Upsert progress
	progress, err := s.progressRepository.Upsert(ctx, userID, missionID, status, watchPositionSeconds)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return progress, nil
}
