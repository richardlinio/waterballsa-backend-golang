package store

import (
	"context"

	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
)

// PayOrderTxParams contains the input parameters for PayOrderTx
type PayOrderTxParams struct {
	OrderID int64
	UserID  int64
}

// PayOrderTxResult contains the result of PayOrderTx
type PayOrderTxResult struct {
	Order      *model.Order
	OrderItems []model.OrderItem
}

// PayOrderTx processes order payment atomically within a transaction.
// It performs the following operations:
// 1. Updates order status to PAID
// 2. Retrieves order items
// 3. Creates user journey ownership records for each item
func (s *Store) PayOrderTx(ctx context.Context, arg PayOrderTxParams) (PayOrderTxResult, error) {
	var result PayOrderTxResult

	err := s.execTx(ctx, func(q *db.Queries) error {
		// Create transaction-aware repositories
		orderRepository := repository.NewOrderRepository(q)
		userJourneyRepository := repository.NewUserJourneyRepository(q)

		// 1. Update order status to PAID
		paidOrder, err := orderRepository.UpdateOrderStatusToPaid(ctx, arg.OrderID)
		if err != nil {
			return err
		}
		result.Order = paidOrder

		// 2. Get order items
		items, err := orderRepository.GetOrderItemsByOrderID(ctx, arg.OrderID)
		if err != nil {
			return err
		}
		result.OrderItems = items

		// 3. Create user journey ownership for each item
		for _, item := range items {
			err := userJourneyRepository.CreateUserJourney(ctx, arg.UserID, item.JourneyID, arg.OrderID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
