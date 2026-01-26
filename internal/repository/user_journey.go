package repository

import (
	"context"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
)

type UserJourneyRepository struct {
	queries db.Querier
}

func NewUserJourneyRepository(queries db.Querier) *UserJourneyRepository {
	return &UserJourneyRepository{queries: queries}
}

// CreateUserJourney creates a user journey ownership record
func (r *UserJourneyRepository) CreateUserJourney(ctx context.Context, userID, journeyID, orderID int64) error {
	return r.queries.CreateUserJourney(ctx, db.CreateUserJourneyParams{
		UserID:    userID,
		JourneyID: journeyID,
		OrderID:   orderID,
	})
}
