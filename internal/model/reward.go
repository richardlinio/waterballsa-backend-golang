package model

import "time"

// Reward represents the domain model for a mission reward
type Reward struct {
	ID          int64
	MissionID   int64
	RewardType  string // EXPERIENCE
	RewardValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
