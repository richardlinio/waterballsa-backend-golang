package model

import "time"

// MissionResource represents the domain model for a mission resource
type MissionResource struct {
	ID              int64
	MissionID       int64
	ResourceType    string // VIDEO, ARTICLE, FORM
	ResourceURL     string
	ResourceContent string
	ContentOrder    int
	DurationSeconds *int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
