package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
)

// missionRepository provides mission data access operations for MissionService
type missionRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Mission, int64, error)
	GetRewardByMissionID(ctx context.Context, missionID int64) (*model.Reward, error)
	ListResourcesByMissionID(ctx context.Context, missionID int64) ([]*model.MissionResource, error)
}

// MissionService implements mission business logic operations.
type MissionService struct {
	missionRepository missionRepository
}

// NewMissionService creates a new MissionService instance.
func NewMissionService(missionRepository *repository.MissionRepository) *MissionService {
	return &MissionService{
		missionRepository: missionRepository,
	}
}

// GetDetail retrieves mission details including reward and resources
// Returns apperror.MissionNotFound if mission doesn't exist
func (s *MissionService) GetDetail(ctx context.Context, missionID int64) (*model.Mission, int64, *model.Reward, []*model.MissionResource, error) {
	// Get mission by ID (also returns journey ID from JOIN)
	mission, journeyID, err := s.missionRepository.GetByID(ctx, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrMissionNotFound) {
			return nil, 0, nil, nil, apperror.MissionNotFound()
		}
		return nil, 0, nil, nil, apperror.DatabaseError(err)
	}

	// Get reward for this mission (no reward means pgx.ErrNoRows, which we handle gracefully)
	reward, err := s.missionRepository.GetRewardByMissionID(ctx, missionID)
	if err != nil {
		// If no reward found, continue with nil reward (not an error)
		if errors.Is(err, pgx.ErrNoRows) {
			reward = nil
		} else {
			return nil, 0, nil, nil, apperror.DatabaseError(err)
		}
	}

	// Get resources for this mission
	resources, err := s.missionRepository.ListResourcesByMissionID(ctx, missionID)
	if err != nil {
		return nil, 0, nil, nil, apperror.DatabaseError(err)
	}

	return mission, journeyID, reward, resources, nil
}
