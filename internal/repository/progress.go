package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
)

var ErrProgressNotFound = errors.New("progress not found")

type ProgressRepository struct {
	queries db.Querier
}

func NewProgressRepository(queries db.Querier) *ProgressRepository {
	return &ProgressRepository{
		queries: queries,
	}
}

// GetByUserAndMission retrieves progress for a specific user and mission
// Uses LEFT JOIN to validate mission existence in a single query
// Returns ErrProgressNotFound if mission doesn't exist
// Returns default progress (UNCOMPLETED, 0 seconds) if mission exists but no progress record
func (r *ProgressRepository) GetByUserAndMission(ctx context.Context, userID, missionID int64) (*model.UserMissionProgress, error) {
	row, err := r.queries.GetUserMissionProgress(ctx, db.GetUserMissionProgressParams{
		Column1: userID,
		ID:      missionID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Mission doesn't exist
			return nil, ErrProgressNotFound
		}
		return nil, err
	}

	return &model.UserMissionProgress{
		ID:                   row.ID,
		UserID:               row.UserID,
		MissionID:            row.MissionID,
		Status:               string(row.Status),
		WatchPositionSeconds: int(row.WatchPositionSeconds),
		CreatedAt:            row.CreatedAt.Time,
		UpdatedAt:            row.UpdatedAt.Time,
	}, nil
}

// Upsert creates or updates progress
func (r *ProgressRepository) Upsert(ctx context.Context, userID, missionID int64, status string, watchPositionSeconds int) (*model.UserMissionProgress, error) {
	//nolint:gosec // watchPositionSeconds is validated in service layer to be within video duration bounds
	row, err := r.queries.UpsertUserMissionProgress(ctx, db.UpsertUserMissionProgressParams{
		UserID:               userID,
		MissionID:            missionID,
		Status:               db.ProgressStatus(status),
		WatchPositionSeconds: int32(watchPositionSeconds),
	})
	if err != nil {
		return nil, err
	}

	return &model.UserMissionProgress{
		ID:                   row.ID,
		UserID:               row.UserID,
		MissionID:            row.MissionID,
		Status:               string(row.Status),
		WatchPositionSeconds: int(row.WatchPositionSeconds),
		CreatedAt:            row.CreatedAt.Time,
		UpdatedAt:            row.UpdatedAt.Time,
	}, nil
}

// ListByUserAndMissions retrieves progress for multiple missions for a specific user
// Returns a map of mission_id -> status for missions that have progress records
// Missions without progress records are not included in the map
func (r *ProgressRepository) ListByUserAndMissions(ctx context.Context, userID int64, missionIDs []int64) (map[int64]string, error) {
	rows, err := r.queries.ListUserMissionProgressByMissions(ctx, db.ListUserMissionProgressByMissionsParams{
		UserID:     userID,
		MissionIds: missionIDs,
	})
	if err != nil {
		return nil, err
	}

	// Convert rows to map
	progressMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		progressMap[row.MissionID] = string(row.Status)
	}

	return progressMap, nil
}
