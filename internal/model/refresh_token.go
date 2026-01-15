package model

import "time"

// RefreshToken represents the domain model for a refresh token
type RefreshToken struct {
	ID        int64
	TokenJTI  string
	UserID    int64
	ExpiresAt time.Time
	RevokedAt *time.Time // nil = valid, non-nil = revoked
	CreatedAt time.Time
}
