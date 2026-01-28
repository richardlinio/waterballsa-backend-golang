package store

import "errors"

var (
	// Auth transaction errors
	ErrTokenUserMismatch = errors.New("token user ID mismatch")

	// Order transaction errors
	ErrJourneyAlreadyPurchased = errors.New("journey already purchased")
	ErrExistingOrderHasNoItems = errors.New("existing order has no items")

	// Progress transaction errors
	ErrMissionAlreadyDelivered = errors.New("mission already delivered")
)
