package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/repository"
	"github.com/richardlinio/waterballsa-backend-golang/internal/store"
)

const (
	// Order number generation constants
	orderNumberTimestampFormat = "2006010215" // YYYYMMDDhh format (10 digits)
	orderNumberRandomBytes     = 3            // 3 bytes = 6 hex chars
	orderNumberRandomCodeLen   = 5            // Take first 5 chars from hex string
)

type orderRepository interface {
	CreateOrder(ctx context.Context, orderNumber string, userID int64, originalPrice, discount, price float64, expiredAt time.Time) (*model.Order, error)
	CreateOrderItem(ctx context.Context, orderID, journeyID int64, quantity int32, originalPrice, discount, price float64) (*model.OrderItem, error)
	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItem, error)
	CheckUserHasPurchasedJourney(ctx context.Context, userID, journeyID int64) (bool, error)
	GetUnpaidOrderByUserAndJourney(ctx context.Context, userID, journeyID int64) (*model.Order, error)
	UpdateOrderStatusToPaid(ctx context.Context, orderID int64) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64, limit, offset int32) ([]model.Order, error)
	CountOrdersByUserID(ctx context.Context, userID int64) (int64, error)
}

type userJourneyRepository interface {
	CreateUserJourney(ctx context.Context, userID, journeyID, orderID int64) error
}

type orderJourneyRepository interface {
	GetByID(ctx context.Context, journeyID int64) (*model.Journey, error)
	GetJourneyTitleByID(ctx context.Context, journeyID int64) (string, error)
}

type orderUserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

// OrderResult holds the complete result of order operations
type OrderResult struct {
	Order         *model.Order
	Items         []model.OrderItem
	JourneyTitles map[int64]string
	Username      string
}

// OrderListResult holds the result of listing orders
type OrderListResult struct {
	Orders        []model.Order
	OrderItems    map[int64][]model.OrderItem // orderID -> items
	JourneyTitles map[int64]string            // journeyID -> title
	Pagination    struct {
		Page  int32
		Limit int32
		Total int64
	}
}

type OrderService struct {
	orderRepository       orderRepository
	journeyRepository     orderJourneyRepository
	userRepository        orderUserRepository
	userJourneyRepository userJourneyRepository
	store                 *store.Store
	transactionTimeout    time.Duration
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	journeyRepository *repository.JourneyRepository,
	userRepository *repository.UserRepository,
	userJourneyRepository *repository.UserJourneyRepository,
	store *store.Store,
	transactionTimeout time.Duration,
) *OrderService {
	return &OrderService{
		orderRepository:       orderRepository,
		journeyRepository:     journeyRepository,
		userRepository:        userRepository,
		userJourneyRepository: userJourneyRepository,
		store:                 store,
		transactionTimeout:    transactionTimeout,
	}
}

// CreateOrder creates a new order or returns existing unpaid order
// Returns (result, isNewOrder, error) where isNewOrder is true for newly created orders
// This method uses a transaction to ensure atomicity and prevent race conditions
func (s *OrderService) CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*OrderResult, bool, error) {
	// 1. Validate that we have exactly one item (MVP constraint)
	if len(req.Items) != 1 {
		return nil, false, apperror.ValidationFailed()
	}

	item := req.Items[0]
	journeyID := item.JourneyID

	// 2. Generate order number (outside TX - no side effects)
	orderNumber, err := s.generateOrderNumber(userID)
	if err != nil {
		return nil, false, apperror.InternalServerError(err)
	}

	// 3. Set parameters for transaction
	discount := 0.0
	expiredAt := time.Now().Add(72 * time.Hour) // 3 days

	// 4. Execute atomic transaction to create or get order
	// This ensures purchase checks, price locking, and order creation are atomic
	// Journey price is fetched and locked inside the transaction
	txCtx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	txResult, err := s.store.CreateOrderTx(txCtx, store.CreateOrderTxParams{
		OrderNumber: orderNumber,
		UserID:      userID,
		JourneyID:   journeyID,
		Quantity:    item.Quantity,
		Discount:    discount,
		ExpiredAt:   expiredAt,
	})
	if err != nil {
		// Check for specific error types using errors.Is()
		if errors.Is(err, repository.ErrJourneyNotFound) {
			return nil, false, apperror.JourneyNotFound()
		}
		if errors.Is(err, store.ErrJourneyAlreadyPurchased) {
			return nil, false, apperror.JourneyAlreadyPurchased()
		}
		return nil, false, apperror.DatabaseError(err)
	}

	// 5. Get user data for response (outside TX - doesn't affect order creation)
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, false, apperror.DatabaseError(err)
	}

	// 6. Build result
	return &OrderResult{
		Order:         txResult.Order,
		Items:         []model.OrderItem{*txResult.OrderItem},
		JourneyTitles: map[int64]string{journeyID: txResult.Journey.Title},
		Username:      user.Username,
	}, txResult.IsNewOrder, nil
}

// GetOrderByID retrieves an order by ID with authorization check
func (s *OrderService) GetOrderByID(ctx context.Context, orderID, userID int64) (*OrderResult, error) {
	order, err := s.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, apperror.OrderNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// Authorization: user can only view their own orders
	if order.UserID != userID {
		return nil, apperror.OrderNotFound() // Return 404 to avoid info leakage
	}

	return s.toOrderResult(ctx, order)
}

// toOrderResult builds a complete order result with items and journey titles
func (s *OrderService) toOrderResult(ctx context.Context, order *model.Order) (*OrderResult, error) {
	// Get order items
	items, err := s.orderRepository.GetOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// Get journey titles
	journeyTitles := make(map[int64]string)
	for _, item := range items {
		title, err := s.journeyRepository.GetJourneyTitleByID(ctx, item.JourneyID)
		if err != nil {
			return nil, apperror.DatabaseError(err)
		}
		journeyTitles[item.JourneyID] = title
	}

	// Get username
	user, err := s.userRepository.GetByID(ctx, order.UserID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &OrderResult{
		Order:         order,
		Items:         items,
		JourneyTitles: journeyTitles,
		Username:      user.Username,
	}, nil
}

// PayOrder processes payment for an order
func (s *OrderService) PayOrder(ctx context.Context, orderID, userID int64) (*model.Order, error) {
	// 1. Get order by ID (pre-transaction validation)
	order, err := s.orderRepository.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, apperror.OrderNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	// 2. Authorization: user can only pay their own orders
	if order.UserID != userID {
		return nil, apperror.OrderNotFound() // Return 404 to avoid info leakage
	}

	// 3. Check order status (pre-transaction validation)
	if order.Status == "PAID" {
		return nil, apperror.OrderAlreadyPaid()
	}
	if order.Status == "EXPIRED" {
		return nil, apperror.OrderExpired()
	}

	// 4. Execute payment transaction through Store
	result, err := s.store.PayOrderTx(ctx, store.PayOrderTxParams{
		OrderID: orderID,
		UserID:  userID,
	})
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return result.Order, nil
}

// GetUserOrders retrieves paginated orders for a user with authorization check
func (s *OrderService) GetUserOrders(ctx context.Context, userID, authenticatedUserID int64, page, limit int32) (*OrderListResult, error) {
	// 1. Authorization: user can only view their own orders
	if userID != authenticatedUserID {
		return nil, apperror.OrderNotFound() // Return 404 to avoid info leakage
	}

	// 2. Validate pagination parameters
	if page < 1 {
		return nil, apperror.ValidationFailed()
	}
	if limit < 1 || limit > 100 {
		return nil, apperror.ValidationFailed()
	}

	// 3. Calculate offset
	offset := (page - 1) * limit

	// 4. Fetch orders with pagination
	orders, err := s.orderRepository.GetOrdersByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// 5. Count total orders for pagination
	total, err := s.orderRepository.CountOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	// 6. Fetch order items and journey titles for each order
	orderItems := make(map[int64][]model.OrderItem)
	journeyTitles := make(map[int64]string)

	for _, order := range orders {
		items, err := s.orderRepository.GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, apperror.DatabaseError(err)
		}
		orderItems[order.ID] = items

		// Fetch journey titles for each item
		for _, item := range items {
			if _, exists := journeyTitles[item.JourneyID]; !exists {
				title, err := s.journeyRepository.GetJourneyTitleByID(ctx, item.JourneyID)
				if err != nil {
					return nil, apperror.DatabaseError(err)
				}
				journeyTitles[item.JourneyID] = title
			}
		}
	}

	// 7. Build result
	result := &OrderListResult{
		Orders:        orders,
		OrderItems:    orderItems,
		JourneyTitles: journeyTitles,
	}
	result.Pagination.Page = page
	result.Pagination.Limit = limit
	result.Pagination.Total = total

	return result, nil
}

// generateOrderNumber generates a unique order number
// Format: {timestamp(10)}{userId}{randomCode(5)}
func (s *OrderService) generateOrderNumber(userID int64) (string, error) {
	// Get current timestamp (10 digits: YYYYMMDDhh)
	timestamp := time.Now().Format(orderNumberTimestampFormat)

	// Generate 5-character random hex code
	randomBytes := make([]byte, orderNumberRandomBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomCode := hex.EncodeToString(randomBytes)[:orderNumberRandomCodeLen]

	return fmt.Sprintf("%s%d%s", timestamp, userID, randomCode), nil
}
