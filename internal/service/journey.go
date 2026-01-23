package service

import (
	"context"

	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
)

// journeyRepository provides journey data access operations for JourneyService
type journeyRepository interface {
	List(ctx context.Context) ([]*model.Journey, error)
}

// JourneyService implements journey business logic operations.
type JourneyService struct {
	journeyRepository journeyRepository
}

// NewJourneyService creates a new JourneyService instance.
func NewJourneyService(journeyRepository *repository.JourneyRepository) *JourneyService {
	return &JourneyService{
		journeyRepository: journeyRepository,
	}
}

// List retrieves all available journeys
func (s *JourneyService) List(ctx context.Context) ([]*model.Journey, error) {
	journeys, err := s.journeyRepository.List(ctx)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return journeys, nil
}
