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
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/service"
	"github.com/richardlinio/waterballsa-backend-golang/internal/util"
)

type orderService interface {
	CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*service.OrderResult, bool, error)
	GetOrderByID(ctx context.Context, orderID, userID int64) (*service.OrderResult, error)
	PayOrder(ctx context.Context, orderID, userID int64) (*model.Order, error)
	GetUserOrders(ctx context.Context, userID, authenticatedUserID int64, page, limit int32) (*service.OrderListResult, error)
}

const (
	defaultPage  = int32(1)
	defaultLimit = int32(20)
	maxLimit     = int32(100)
)

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

// PayOrder handles POST /orders/:orderId/action/pay
func (h *OrderHandler) PayOrder(c *gin.Context) {
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

	paidOrder, err := h.orderService.PayOrder(ctx, orderID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	message := "付款完成"
	response := dto.ToPayOrderResponse(paidOrder, message)
	c.JSON(http.StatusOK, response)
}

// GetUserOrders handles GET /users/:userId/orders
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	// Parse userId from path parameter
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
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

	// Parse pagination query parameters with defaults
	page := defaultPage
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err == nil && pageInt > 0 {
			page = int32(pageInt)
		}
	}

	limit := defaultLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err == nil && limitInt > 0 && limitInt <= int64(maxLimit) {
			limit = int32(limitInt)
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Call service with authorization check
	result, err := h.orderService.GetUserOrders(ctx, userID, authenticatedUser.ID, page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Convert to DTO response
	response := dto.ToOrderListResponse(
		result.Orders,
		result.OrderItems,
		result.JourneyTitles,
		result.Pagination.Page,
		result.Pagination.Limit,
		result.Pagination.Total,
	)
	c.JSON(http.StatusOK, response)
}
