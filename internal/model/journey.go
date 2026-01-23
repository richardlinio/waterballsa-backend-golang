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

// Mission represents the domain model for a mission
type Mission struct {
	ID          int64
	ChapterID   int64
	Title       string
	Type        string // VIDEO, ARTICLE, QUESTIONNAIRE
	AccessLevel string // PUBLIC, AUTHENTICATED, PURCHASED
	OrderIndex  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
