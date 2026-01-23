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
