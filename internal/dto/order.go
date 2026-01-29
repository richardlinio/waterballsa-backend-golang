package dto

import (
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
)

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	JourneyID int64 `json:"journeyId" binding:"required,gt=0"`
	Quantity  int32 `json:"quantity" binding:"required,min=1"`
}

type OrderItemDTO struct {
	JourneyID     int64   `json:"journeyId"`
	JourneyTitle  string  `json:"journeyTitle"`
	Quantity      int32   `json:"quantity"`
	OriginalPrice float64 `json:"originalPrice"`
	Discount      float64 `json:"discount"`
	Price         float64 `json:"price"`
}

type OrderResponse struct {
	ID            int64          `json:"id"`
	OrderNumber   string         `json:"orderNumber"`
	UserID        int64          `json:"userId"`
	Username      string         `json:"username"`
	Status        string         `json:"status"`
	OriginalPrice float64        `json:"originalPrice"`
	Discount      float64        `json:"discount"`
	Price         float64        `json:"price"`
	Items         []OrderItemDTO `json:"items"`
	CreatedAt     int64          `json:"createdAt"` // Unix milliseconds
	ExpiredAt     *int64         `json:"expiredAt"` // Unix milliseconds, nullable
	PaidAt        *int64         `json:"paidAt"`    // Unix milliseconds, nullable
}

// ToOrderResponse converts order data to OrderResponse DTO
func ToOrderResponse(
	order *model.Order,
	items []model.OrderItem,
	journeyTitles map[int64]string,
	username string,
) OrderResponse {
	// Convert items
	itemDTOs := make([]OrderItemDTO, 0, len(items))
	for _, item := range items {
		itemDTOs = append(itemDTOs, OrderItemDTO{
			JourneyID:     item.JourneyID,
			JourneyTitle:  journeyTitles[item.JourneyID],
			Quantity:      item.Quantity,
			OriginalPrice: item.OriginalPrice,
			Discount:      item.Discount,
			Price:         item.Price,
		})
	}

	// Convert timestamps to Unix milliseconds
	createdAt := order.CreatedAt.UnixMilli()

	var expiredAt *int64
	if order.ExpiredAt != nil {
		ts := order.ExpiredAt.UnixMilli()
		expiredAt = &ts
	}

	var paidAt *int64
	if order.PaidAt != nil {
		ts := order.PaidAt.UnixMilli()
		paidAt = &ts
	}

	return OrderResponse{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		UserID:        order.UserID,
		Username:      username,
		Status:        order.Status,
		OriginalPrice: order.OriginalPrice,
		Discount:      order.Discount,
		Price:         order.Price,
		Items:         itemDTOs,
		CreatedAt:     createdAt,
		ExpiredAt:     expiredAt,
		PaidAt:        paidAt,
	}
}

type PayOrderResponse struct {
	ID          int64   `json:"id"`
	OrderNumber string  `json:"orderNumber"`
	Status      string  `json:"status"`
	Price       float64 `json:"price"`
	PaidAt      int64   `json:"paidAt"`
	Message     string  `json:"message"`
}

// ToPayOrderResponse converts paid order data to PayOrderResponse DTO
func ToPayOrderResponse(order *model.Order, message string) PayOrderResponse {
	return PayOrderResponse{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		Price:       order.Price,
		PaidAt:      order.PaidAt.UnixMilli(),
		Message:     message,
	}
}

// OrderItemSummary is a simplified order item for list view
type OrderItemSummary struct {
	JourneyID    int64  `json:"journeyId"`
	JourneyTitle string `json:"journeyTitle"`
}

// OrderSummary is a simplified order for list view
type OrderSummary struct {
	ID          int64              `json:"id"`
	OrderNumber string             `json:"orderNumber"`
	Status      string             `json:"status"`
	Price       float64            `json:"price"`
	Items       []OrderItemSummary `json:"items"`
	CreatedAt   int64              `json:"createdAt"` // Unix milliseconds
	ExpiredAt   *int64             `json:"expiredAt"` // Unix milliseconds, nullable
	PaidAt      *int64             `json:"paidAt"`    // Unix milliseconds, nullable
}

// Pagination holds pagination metadata
type Pagination struct {
	Page  int32 `json:"page"`  // Current page (1-indexed)
	Limit int32 `json:"limit"` // Items per page
	Total int64 `json:"total"` // Total items count
}

// OrderListResponse is the response for listing orders
type OrderListResponse struct {
	Orders     []OrderSummary `json:"orders"`
	Pagination Pagination     `json:"pagination"`
}

// ToOrderListResponse converts order list data to OrderListResponse DTO
func ToOrderListResponse(
	orders []model.Order,
	orderItems map[int64][]model.OrderItem,
	journeyTitles map[int64]string,
	page, limit int32,
	total int64,
) OrderListResponse {
	orderSummaries := make([]OrderSummary, 0, len(orders))

	for _, order := range orders {
		// Convert items for this order
		items := orderItems[order.ID]
		itemSummaries := make([]OrderItemSummary, 0, len(items))
		for _, item := range items {
			itemSummaries = append(itemSummaries, OrderItemSummary{
				JourneyID:    item.JourneyID,
				JourneyTitle: journeyTitles[item.JourneyID],
			})
		}

		// Convert timestamps to Unix milliseconds
		createdAt := order.CreatedAt.UnixMilli()

		var expiredAt *int64
		if order.ExpiredAt != nil {
			ts := order.ExpiredAt.UnixMilli()
			expiredAt = &ts
		}

		var paidAt *int64
		if order.PaidAt != nil {
			ts := order.PaidAt.UnixMilli()
			paidAt = &ts
		}

		orderSummaries = append(orderSummaries, OrderSummary{
			ID:          order.ID,
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			Price:       order.Price,
			Items:       itemSummaries,
			CreatedAt:   createdAt,
			ExpiredAt:   expiredAt,
			PaidAt:      paidAt,
		})
	}

	return OrderListResponse{
		Orders: orderSummaries,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}
