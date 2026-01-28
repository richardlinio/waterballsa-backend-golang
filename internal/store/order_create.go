package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

// CreateOrderTxParams contains the input parameters for CreateOrderTx
type CreateOrderTxParams struct {
	OrderNumber string
	UserID      int64
	JourneyID   int64
	Quantity    int32
	Discount    float64
	ExpiredAt   time.Time
}

// CreateOrderTxResult contains the result of CreateOrderTx
type CreateOrderTxResult struct {
	Order      *model.Order
	OrderItem  *model.OrderItem
	Journey    *model.Journey
	IsNewOrder bool // true if created new order, false if returning existing unpaid order
}

// CreateOrderTx creates a new order atomically within a transaction.
// It performs the following operations:
// 1. Checks if user has already purchased the journey (prevents duplicate purchases)
// 2. Checks for existing unpaid order for this user+journey
// 3. If existing unpaid order found:
//   - Validates it's not expired
//   - Returns existing order with IsNewOrder=false
//
// 4. If no existing unpaid order:
//   - Gets journey by ID to lock price
//   - Creates new order
//   - Creates order item
//   - Returns new order with IsNewOrder=true
//
// Lock ordering follows: journeys → orders → order_items → user_journeys
// This prevents deadlocks with PayOrderTx which uses the same order.
//
//nolint:gocyclo // Complex transaction logic requires multiple validation branches
func (s *Store) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (CreateOrderTxResult, error) {
	var result CreateOrderTxResult

	err := s.execTx(ctx, func(q *db.Queries) error {
		// Create transaction-aware repositories
		orderRepository := repository.NewOrderRepository(q)
		journeyRepository := repository.NewJourneyRepository(q)

		// 1. Check if user has already purchased this journey
		// This prevents duplicate purchases via TOCTOU race conditions
		hasPurchased, err := orderRepository.CheckUserHasPurchasedJourney(ctx, arg.UserID, arg.JourneyID)
		if err != nil {
			return fmt.Errorf("failed to check purchase status: %w", err)
		}
		if hasPurchased {
			return fmt.Errorf("journey already purchased")
		}

		// 2. Check for existing unpaid order
		existingOrder, err := orderRepository.GetUnpaidOrderByUserAndJourney(ctx, arg.UserID, arg.JourneyID)
		if err != nil && !errors.Is(err, repository.ErrNoUnpaidOrderFound) {
			return fmt.Errorf("failed to check unpaid orders: %w", err)
		}

		// 3. If existing unpaid order found, validate and return it
		if err == nil && existingOrder != nil {
			// Validate order is not expired
			if existingOrder.ExpiredAt != nil && existingOrder.ExpiredAt.Before(time.Now()) {
				// Order expired, continue to create new order instead
			} else {
				// Return existing valid unpaid order
				items, err := orderRepository.GetOrderItemsByOrderID(ctx, existingOrder.ID)
				if err != nil {
					return fmt.Errorf("failed to get order items: %w", err)
				}
				if len(items) == 0 {
					return fmt.Errorf("existing order has no items")
				}

				// Get journey for response
				journey, err := journeyRepository.GetByID(ctx, arg.JourneyID)
				if err != nil {
					return fmt.Errorf("failed to get journey: %w", err)
				}

				result.Order = existingOrder
				result.OrderItem = &items[0]
				result.Journey = journey
				result.IsNewOrder = false
				return nil
			}
		}

		// 4. No existing valid unpaid order - create new order

		// Get journey by ID to lock price (prevents price changes during order creation)
		journey, err := journeyRepository.GetByID(ctx, arg.JourneyID)
		if err != nil {
			if errors.Is(err, repository.ErrJourneyNotFound) {
				return fmt.Errorf("journey not found")
			}
			return fmt.Errorf("failed to get journey: %w", err)
		}

		// Calculate prices from journey
		originalPrice := journey.Price
		price := originalPrice - arg.Discount

		// Create order
		order, err := orderRepository.CreateOrder(
			ctx,
			arg.OrderNumber,
			arg.UserID,
			originalPrice,
			arg.Discount,
			price,
			arg.ExpiredAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create order item
		orderItem, err := orderRepository.CreateOrderItem(
			ctx,
			order.ID,
			arg.JourneyID,
			arg.Quantity,
			originalPrice,
			arg.Discount,
			price,
		)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}

		result.Order = order
		result.OrderItem = orderItem
		result.Journey = journey
		result.IsNewOrder = true
		return nil
	})
	if err != nil {
		return CreateOrderTxResult{}, fmt.Errorf("create order transaction failed: %w", err)
	}

	return result, nil
}
