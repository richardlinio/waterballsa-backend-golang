package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
)

// ErrJourneyNotFound is returned when a journey is not found in the database
var ErrJourneyNotFound = errors.New("journey not found")

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

// GetByID retrieves a journey by its ID
// Returns ErrJourneyNotFound if journey is not found
func (r *JourneyRepository) GetByID(ctx context.Context, id int64) (*model.Journey, error) {
	row, err := r.queries.GetJourneyByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJourneyNotFound
		}
		return nil, err
	}

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

	return journey, nil
}

// ListChaptersByJourneyID retrieves all chapters for a journey, ordered by order_index
func (r *JourneyRepository) ListChaptersByJourneyID(ctx context.Context, journeyID int64) ([]*model.Chapter, error) {
	rows, err := r.queries.ListChaptersByJourneyID(ctx, journeyID)
	if err != nil {
		return nil, err
	}

	chapters := make([]*model.Chapter, 0, len(rows))
	for _, row := range rows {
		chapter := &model.Chapter{
			ID:         row.ID,
			JourneyID:  row.JourneyID,
			Title:      row.Title,
			OrderIndex: int(row.OrderIndex),
			CreatedAt:  row.CreatedAt.Time,
			UpdatedAt:  row.UpdatedAt.Time,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

// ListMissionsByChapterIDs retrieves all missions for the given chapter IDs
// Results are ordered by chapter_id and order_index
func (r *JourneyRepository) ListMissionsByChapterIDs(ctx context.Context, chapterIDs []int64) ([]*model.Mission, error) {
	rows, err := r.queries.ListMissionsByChapterIDs(ctx, chapterIDs)
	if err != nil {
		return nil, err
	}

	missions := make([]*model.Mission, 0, len(rows))
	for _, row := range rows {
		mission := &model.Mission{
			ID:          row.ID,
			ChapterID:   row.ChapterID,
			Title:       row.Title,
			Type:        string(row.Type),
			AccessLevel: string(row.AccessLevel),
			OrderIndex:  int(row.OrderIndex),
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
		missions = append(missions, mission)
	}

	return missions, nil
}

// GetJourneyTitleByID retrieves journey title by ID
func (r *JourneyRepository) GetJourneyTitleByID(ctx context.Context, journeyID int64) (string, error) {
	title, err := r.queries.GetJourneyTitleByID(ctx, journeyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrJourneyNotFound
		}
		return "", err
	}
	return title, nil
}

// UpdateJourneyPrice updates the price of a journey
func (r *JourneyRepository) UpdateJourneyPrice(ctx context.Context, journeyID int64, price float64) error {
	var priceNumeric pgtype.Numeric
	if err := priceNumeric.Scan(fmt.Sprintf("%.2f", price)); err != nil {
		return err
	}

	err := r.queries.UpdateJourneyPrice(ctx, db.UpdateJourneyPriceParams{
		ID:    journeyID,
		Price: priceNumeric,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJourneyNotFound
		}
		return err
	}
	return nil
}
