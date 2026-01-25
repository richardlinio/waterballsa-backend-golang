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

type OrderItemDTO struct {
	JourneyID     int64   `json:"journeyId"`
	JourneyTitle  string  `json:"journeyTitle"`
	Quantity      int32   `json:"quantity"`
	OriginalPrice float64 `json:"originalPrice"`
	Discount      float64 `json:"discount"`
	Price         float64 `json:"price"`
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
