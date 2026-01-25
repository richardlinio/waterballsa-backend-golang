package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/service"
	"github.com/richardlinio/waterballsa-backend-golang/internal/util"
)

type orderService interface {
	CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*service.OrderResult, bool, error)
	GetOrderByID(ctx context.Context, orderID, userID int64) (*service.OrderResult, error)
}

type OrderHandler struct {
	orderService   orderService
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewOrderHandler(
	orderService *service.OrderService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *OrderHandler {
	return &OrderHandler{
		orderService:   orderService,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	// Get authenticated user from JWT
	authenticatedUser := util.GetAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}
	userID := authenticatedUser.ID

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	result, isNewOrder, err := h.orderService.CreateOrder(ctx, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Determine status code based on whether this is a new order or existing order
	statusCode := http.StatusOK
	if isNewOrder {
		statusCode = http.StatusCreated
	}

	response := dto.ToOrderResponse(result.Order, result.Items, result.JourneyTitles, result.Username)
	c.JSON(statusCode, response)
}

// GetOrderDetail handles GET /orders/:orderId
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	orderIDStr := c.Param("orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		_ = c.Error(apperror.ValidationFailed())
		return
	}

	// Get authenticated user from JWT
	authenticatedUser := util.GetAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}
	userID := authenticatedUser.ID

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	result, err := h.orderService.GetOrderByID(ctx, orderID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToOrderResponse(result.Order, result.Items, result.JourneyTitles, result.Username)
	c.JSON(http.StatusOK, response)
}
