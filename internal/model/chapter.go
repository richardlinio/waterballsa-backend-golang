package model

import "time"

// Chapter represents the domain model for a chapter
type Chapter struct {
	ID         int64
	JourneyID  int64
	Title      string
	OrderIndex int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
