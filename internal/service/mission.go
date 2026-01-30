package service

import (
	"context"
	"errors"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
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
	orderRepository   orderRepository
}

// NewMissionService creates a new MissionService instance.
func NewMissionService(missionRepository *repository.MissionRepository, orderRepository *repository.OrderRepository) *MissionService {
	return &MissionService{
		missionRepository: missionRepository,
		orderRepository:   orderRepository,
	}
}

// GetDetail retrieves mission details including reward and resources
// Performs access control checks based on mission access level
// Returns apperror.MissionNotFound if mission doesn't exist
// Returns apperror.Unauthorized if authentication required but user not provided
// Returns apperror.Forbidden if user hasn't purchased required journey
func (s *MissionService) GetDetail(ctx context.Context, missionID int64, userID *int64) (*model.MissionDetail, error) {
	// Get mission by ID (also returns journey ID from JOIN)
	mission, journeyID, err := s.missionRepository.GetByID(ctx, missionID)
	if err != nil {
		if errors.Is(err, repository.ErrMissionNotFound) {
			return nil, apperror.MissionNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Access control: check authentication requirement
	if mission.AccessLevel == "AUTHENTICATED" || mission.AccessLevel == "PURCHASED" {
		if userID == nil {
			return nil, apperror.Unauthorized()
		}
	}

	// Access control: check purchase requirement
	if mission.AccessLevel == "PURCHASED" {
		hasPurchased, err := s.orderRepository.CheckUserHasPurchasedJourney(ctx, *userID, journeyID)
		if err != nil {
			return nil, apperror.DatabaseError(err)
		}
		if !hasPurchased {
			return nil, apperror.Forbidden()
		}
	}

	// Get reward for this mission (no reward means repository.ErrRewardNotFound, which we handle gracefully)
	reward, err := s.missionRepository.GetRewardByMissionID(ctx, missionID)
	if err != nil {
		// If no reward found, continue with nil reward (not an error)
		if errors.Is(err, repository.ErrRewardNotFound) {
			reward = nil
		} else {
			return nil, apperror.DatabaseError(err)
		}
	}

	// Get resources for this mission
	resources, err := s.missionRepository.ListResourcesByMissionID(ctx, missionID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &model.MissionDetail{
		Mission:   mission,
		JourneyID: journeyID,
		Reward:    reward,
		Resources: resources,
	}, nil
}
