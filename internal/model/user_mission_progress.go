package model

import "time"

// UserMissionProgress represents user's progress on a mission
type UserMissionProgress struct {
	ID                   int64
	UserID               int64
	MissionID            int64
	Status               string // UNCOMPLETED, COMPLETED, DELIVERED
	WatchPositionSeconds int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}
