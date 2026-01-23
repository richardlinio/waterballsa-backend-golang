package model

import "time"

// Mission represents the domain model for a mission
type Mission struct {
	ID          int64
	ChapterID   int64
	Title       string
	Description string
	Type        string // VIDEO, ARTICLE, QUESTIONNAIRE
	AccessLevel string // PUBLIC, AUTHENTICATED, PURCHASED
	OrderIndex  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// MissionDetail represents the complete view of a mission with its related data
// This is an aggregate that combines mission, reward, and resources for read operations
type MissionDetail struct {
	Mission   *Mission
	JourneyID int64
	Reward    *Reward
	Resources []*MissionResource
}
