package repository

import (
	"context"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
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

// GetUserJourneys retrieves all purchased journeys for a user
func (r *UserJourneyRepository) GetUserJourneys(ctx context.Context, userID int64) ([]dto.UserJourneyItem, error) {
	rows, err := r.queries.GetUserJourneysByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	journeys := make([]dto.UserJourneyItem, 0, len(rows))
	for _, row := range rows {
		journeys = append(journeys, dto.UserJourneyItem{
			JourneyID:     row.JourneyID,
			JourneyTitle:  row.JourneyTitle,
			JourneySlug:   row.JourneySlug,
			CoverImageURL: row.CoverImageUrl.String,
			TeacherName:   row.TeacherName,
			PurchasedAt:   row.PurchasedAt.Time.UnixMilli(),
			OrderNumber:   row.OrderNumber,
		})
	}

	return journeys, nil
}
