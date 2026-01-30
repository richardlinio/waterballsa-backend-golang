package service

import (
	"context"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

type userJourneyQueryRepository interface {
	GetUserJourneys(ctx context.Context, userID int64) ([]dto.UserJourneyItem, error)
}

type UserJourneyService struct {
	userJourneyRepository userJourneyQueryRepository
}

func NewUserJourneyService(
	userJourneyRepository *repository.UserJourneyRepository,
) *UserJourneyService {
	return &UserJourneyService{
		userJourneyRepository: userJourneyRepository,
	}
}

// GetUserJourneys retrieves all purchased journeys for a user with authorization check
func (s *UserJourneyService) GetUserJourneys(ctx context.Context, userID, authenticatedUserID int64) ([]dto.UserJourneyItem, error) {
	// 1. Authorization: user can only view their own journeys
	if userID != authenticatedUserID {
		return nil, apperror.OrderNotFound() // Return 404 to avoid info leakage
	}

	// 2. Fetch user's purchased journeys from repository
	journeys, err := s.userJourneyRepository.GetUserJourneys(ctx, userID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return journeys, nil
}
