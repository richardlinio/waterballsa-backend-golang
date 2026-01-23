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

// JourneyDetail represents the complete view of a journey with its related data
// This is an aggregate that combines journey, chapters, and missions for read operations
type JourneyDetail struct {
	Journey           *Journey
	Chapters          []*Chapter
	MissionsByChapter map[int64][]*Mission
}
