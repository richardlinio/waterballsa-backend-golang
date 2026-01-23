package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// ErrMissionNotFound is returned when a mission is not found in the database
var ErrMissionNotFound = errors.New("mission not found")

// MissionRepository implements mission data access operations using sqlc generated queries.
type MissionRepository struct {
	queries db.Querier
}

// NewMissionRepository creates a new instance of MissionRepository
func NewMissionRepository(queries db.Querier) *MissionRepository {
	return &MissionRepository{
		queries: queries,
	}
}

// GetByID retrieves a mission by its ID along with the journey ID
// Returns ErrMissionNotFound if mission is not found
func (r *MissionRepository) GetByID(ctx context.Context, id int64) (*model.Mission, int64, error) {
	row, err := r.queries.GetMissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, ErrMissionNotFound
		}
		return nil, 0, err
	}

	mission := &model.Mission{
		ID:          row.ID,
		ChapterID:   row.ChapterID,
		Title:       row.Title,
		Description: row.Description.String,
		Type:        string(row.Type),
		AccessLevel: string(row.AccessLevel),
		OrderIndex:  int(row.OrderIndex),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return mission, row.JourneyID, nil
}

// GetRewardByMissionID retrieves the reward for a mission
// Returns database error if the query fails
func (r *MissionRepository) GetRewardByMissionID(ctx context.Context, missionID int64) (*model.Reward, error) {
	row, err := r.queries.GetRewardByMissionID(ctx, missionID)
	if err != nil {
		return nil, err
	}

	reward := &model.Reward{
		ID:          row.ID,
		MissionID:   row.MissionID,
		RewardType:  string(row.RewardType),
		RewardValue: int(row.RewardValue),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return reward, nil
}

// ListResourcesByMissionID retrieves all resources for a mission, ordered by content_order
func (r *MissionRepository) ListResourcesByMissionID(ctx context.Context, missionID int64) ([]*model.MissionResource, error) {
	rows, err := r.queries.ListResourcesByMissionID(ctx, missionID)
	if err != nil {
		return nil, err
	}

	resources := make([]*model.MissionResource, 0, len(rows))
	for _, row := range rows {
		var durationSeconds *int
		if row.DurationSeconds.Valid {
			duration := int(row.DurationSeconds.Int32)
			durationSeconds = &duration
		}

		resource := &model.MissionResource{
			ID:              row.ID,
			MissionID:       row.MissionID,
			ResourceType:    string(row.ResourceType),
			ResourceURL:     row.ResourceUrl.String,
			ResourceContent: row.ResourceContent.String,
			ContentOrder:    int(row.ContentOrder),
			DurationSeconds: durationSeconds,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
		}
		resources = append(resources, resource)
	}

	return resources, nil
}
