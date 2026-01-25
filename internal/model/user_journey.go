package model

import "time"

type UserJourney struct {
	ID          int64
	UserID      int64
	JourneyID   int64
	OrderID     int64
	PurchasedAt time.Time
	CreatedAt   time.Time
}
