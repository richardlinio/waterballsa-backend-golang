package model

import "time"

// Journey represents the domain model for a journey
type Journey struct {
	ID            int64
	Title         string
	Slug          string
	Description   string
	CoverImageURL string
	TeacherName   string
	Price         float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
