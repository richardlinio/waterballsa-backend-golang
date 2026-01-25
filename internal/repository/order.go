package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardlinio/waterballsa-backend-golang/internal/db"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrNoUnpaidOrderFound = errors.New("no unpaid order found")
)

type OrderRepository struct {
	queries db.Querier
}

func NewOrderRepository(queries db.Querier) *OrderRepository {
	return &OrderRepository{queries: queries}
}

// CreateOrder creates a new order and returns the created order
func (r *OrderRepository) CreateOrder(ctx context.Context, orderNumber string, userID int64, originalPrice, discount, price float64, expiredAt time.Time) (*model.Order, error) {
	var originalPriceNumeric, discountNumeric, priceNumeric pgtype.Numeric
	if err := originalPriceNumeric.Scan(fmt.Sprintf("%.2f", originalPrice)); err != nil {
		return nil, err
	}
	if err := discountNumeric.Scan(fmt.Sprintf("%.2f", discount)); err != nil {
		return nil, err
	}
	if err := priceNumeric.Scan(fmt.Sprintf("%.2f", price)); err != nil {
		return nil, err
	}

	row, err := r.queries.CreateOrder(ctx, db.CreateOrderParams{
		OrderNumber:   orderNumber,
		UserID:        userID,
		Status:        "UNPAID",
		OriginalPrice: originalPriceNumeric,
		Discount:      discountNumeric,
		Price:         priceNumeric,
		ExpiredAt:     pgtype.Timestamp{Time: expiredAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Convert numeric values
	originalPriceFloat, err := row.OriginalPrice.Float64Value()
	if err != nil {
		return nil, err
	}
	discountFloat, err := row.Discount.Float64Value()
	if err != nil {
		return nil, err
	}
	priceFloat, err := row.Price.Float64Value()
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		ID:            row.ID,
		OrderNumber:   row.OrderNumber,
		UserID:        row.UserID,
		Status:        string(row.Status),
		OriginalPrice: originalPriceFloat.Float64,
		Discount:      discountFloat.Float64,
		Price:         priceFloat.Float64,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}

	if row.ExpiredAt.Valid {
		order.ExpiredAt = &row.ExpiredAt.Time
	}
	if row.PaidAt.Valid {
		order.PaidAt = &row.PaidAt.Time
	}

	return order, nil
}

// CreateOrderItem creates a new order item
func (r *OrderRepository) CreateOrderItem(ctx context.Context, orderID, journeyID int64, quantity int32, originalPrice, discount, price float64) (*model.OrderItem, error) {
	var originalPriceNumeric, discountNumeric, priceNumeric pgtype.Numeric
	if err := originalPriceNumeric.Scan(fmt.Sprintf("%.2f", originalPrice)); err != nil {
		return nil, err
	}
	if err := discountNumeric.Scan(fmt.Sprintf("%.2f", discount)); err != nil {
		return nil, err
	}
	if err := priceNumeric.Scan(fmt.Sprintf("%.2f", price)); err != nil {
		return nil, err
	}

	row, err := r.queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:       orderID,
		JourneyID:     journeyID,
		Quantity:      quantity,
		OriginalPrice: originalPriceNumeric,
		Discount:      discountNumeric,
		Price:         priceNumeric,
	})
	if err != nil {
		return nil, err
	}

	// Convert numeric values
	originalPriceFloat, err := row.OriginalPrice.Float64Value()
	if err != nil {
		return nil, err
	}
	discountFloat, err := row.Discount.Float64Value()
	if err != nil {
		return nil, err
	}
	priceFloat, err := row.Price.Float64Value()
	if err != nil {
		return nil, err
	}

	return &model.OrderItem{
		ID:            row.ID,
		OrderID:       row.OrderID,
		JourneyID:     row.JourneyID,
		Quantity:      row.Quantity,
		OriginalPrice: originalPriceFloat.Float64,
		Discount:      discountFloat.Float64,
		Price:         priceFloat.Float64,
		CreatedAt:     row.CreatedAt.Time,
	}, nil
}

// GetOrderByID retrieves an order by ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	originalPriceFloat, err := row.OriginalPrice.Float64Value()
	if err != nil {
		return nil, err
	}
	discountFloat, err := row.Discount.Float64Value()
	if err != nil {
		return nil, err
	}
	priceFloat, err := row.Price.Float64Value()
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		ID:            row.ID,
		OrderNumber:   row.OrderNumber,
		UserID:        row.UserID,
		Status:        string(row.Status),
		OriginalPrice: originalPriceFloat.Float64,
		Discount:      discountFloat.Float64,
		Price:         priceFloat.Float64,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}

	if row.ExpiredAt.Valid {
		order.ExpiredAt = &row.ExpiredAt.Time
	}
	if row.PaidAt.Valid {
		order.PaidAt = &row.PaidAt.Time
	}

	return order, nil
}

// GetOrderItemsByOrderID retrieves all order items for an order
func (r *OrderRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItem, error) {
	rows, err := r.queries.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items := make([]model.OrderItem, 0, len(rows))
	for _, row := range rows {
		originalPriceFloat, err := row.OriginalPrice.Float64Value()
		if err != nil {
			return nil, err
		}
		discountFloat, err := row.Discount.Float64Value()
		if err != nil {
			return nil, err
		}
		priceFloat, err := row.Price.Float64Value()
		if err != nil {
			return nil, err
		}

		items = append(items, model.OrderItem{
			ID:            row.ID,
			OrderID:       row.OrderID,
			JourneyID:     row.JourneyID,
			Quantity:      row.Quantity,
			OriginalPrice: originalPriceFloat.Float64,
			Discount:      discountFloat.Float64,
			Price:         priceFloat.Float64,
			CreatedAt:     row.CreatedAt.Time,
		})
	}

	return items, nil
}

// CheckUserHasPurchasedJourney checks if user has already purchased the journey
func (r *OrderRepository) CheckUserHasPurchasedJourney(ctx context.Context, userID, journeyID int64) (bool, error) {
	hasPurchased, err := r.queries.CheckUserHasPurchasedJourney(ctx, db.CheckUserHasPurchasedJourneyParams{
		UserID:    userID,
		JourneyID: journeyID,
	})
	if err != nil {
		return false, err
	}
	return hasPurchased, nil
}

// GetUnpaidOrderByUserAndJourney retrieves an unpaid order for a user and journey
func (r *OrderRepository) GetUnpaidOrderByUserAndJourney(ctx context.Context, userID, journeyID int64) (*model.Order, error) {
	row, err := r.queries.GetUnpaidOrderByUserAndJourney(ctx, db.GetUnpaidOrderByUserAndJourneyParams{
		UserID:    userID,
		JourneyID: journeyID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoUnpaidOrderFound
		}
		return nil, err
	}

	originalPriceFloat, err := row.OriginalPrice.Float64Value()
	if err != nil {
		return nil, err
	}
	discountFloat, err := row.Discount.Float64Value()
	if err != nil {
		return nil, err
	}
	priceFloat, err := row.Price.Float64Value()
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		ID:            row.ID,
		OrderNumber:   row.OrderNumber,
		UserID:        row.UserID,
		Status:        string(row.Status),
		OriginalPrice: originalPriceFloat.Float64,
		Discount:      discountFloat.Float64,
		Price:         priceFloat.Float64,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}

	if row.ExpiredAt.Valid {
		order.ExpiredAt = &row.ExpiredAt.Time
	}
	if row.PaidAt.Valid {
		order.PaidAt = &row.PaidAt.Time
	}

	return order, nil
}
