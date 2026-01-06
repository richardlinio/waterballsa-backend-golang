package model

import "time"

// User represents the domain model for a user
type User struct {
	ID               int64
	Username         string
	PasswordHash     string
	Role             string
	ExperiencePoints int32
	Level            int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
