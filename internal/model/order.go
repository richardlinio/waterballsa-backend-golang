package model

import "time"

type Order struct {
	ID            int64
	OrderNumber   string
	UserID        int64
	Status        string // UNPAID, PAID, EXPIRED
	OriginalPrice float64
	Discount      float64
	Price         float64
	CreatedAt     time.Time
	ExpiredAt     *time.Time
	PaidAt        *time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID            int64
	OrderID       int64
	JourneyID     int64
	Quantity      int32
	OriginalPrice float64
	Discount      float64
	Price         float64
	CreatedAt     time.Time
}
