package dto

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
