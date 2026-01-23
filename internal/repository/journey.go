package repository

import (
	"context"

	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// JourneyRepository implements journey data access operations using sqlc generated queries.
type JourneyRepository struct {
	queries db.Querier
}

// NewJourneyRepository creates a new instance of JourneyRepository
func NewJourneyRepository(queries db.Querier) *JourneyRepository {
	return &JourneyRepository{
		queries: queries,
	}
}

// List retrieves all non-deleted journeys ordered by creation date
func (r *JourneyRepository) List(ctx context.Context) ([]*model.Journey, error) {
	rows, err := r.queries.ListJourneys(ctx)
	if err != nil {
		return nil, err
	}

	// Convert sqlc-generated rows to domain models
	journeys := make([]*model.Journey, 0, len(rows))
	for _, row := range rows {
		// Convert pgtype.Numeric to float64
		priceFloat, err := row.Price.Float64Value()
		if err != nil {
			return nil, err
		}

		journey := &model.Journey{
			ID:            row.ID,
			Title:         row.Title,
			Slug:          row.Slug,
			Description:   row.Description.String,
			CoverImageURL: row.CoverImageUrl.String,
			TeacherName:   row.TeacherName,
			Price:         priceFloat.Float64,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
		}
		journeys = append(journeys, journey)
	}

	return journeys, nil
}
