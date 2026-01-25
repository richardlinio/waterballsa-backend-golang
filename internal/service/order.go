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
)

const (
	// Order number generation constants
	orderNumberTimestampFormat = "2006010215" // YYYYMMDDhh format (10 digits)
	orderNumberRandomBytes     = 3            // 3 bytes = 6 hex chars
	orderNumberRandomCodeLen   = 5            // Take first 5 chars from hex string
)

type orderRepositoryForOrder interface {
	CreateOrder(ctx context.Context, orderNumber string, userID int64, originalPrice, discount, price float64, expiredAt time.Time) (*model.Order, error)
	CreateOrderItem(ctx context.Context, orderID, journeyID int64, quantity int32, originalPrice, discount, price float64) (*model.OrderItem, error)
	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]model.OrderItem, error)
	CheckUserHasPurchasedJourney(ctx context.Context, userID, journeyID int64) (bool, error)
	GetUnpaidOrderByUserAndJourney(ctx context.Context, userID, journeyID int64) (*model.Order, error)
}

type journeyRepositoryForOrder interface {
	GetByID(ctx context.Context, journeyID int64) (*model.Journey, error)
	GetJourneyTitleByID(ctx context.Context, journeyID int64) (string, error)
}

type userRepositoryForOrder interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

// OrderResult holds the complete result of order operations
type OrderResult struct {
	Order         *model.Order
	Items         []model.OrderItem
	JourneyTitles map[int64]string
	Username      string
}

type OrderService struct {
	orderRepository   orderRepositoryForOrder
	journeyRepository journeyRepositoryForOrder
	userRepository    userRepositoryForOrder
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	journeyRepository *repository.JourneyRepository,
	userRepository *repository.UserRepository,
) *OrderService {
	return &OrderService{
		orderRepository:   orderRepository,
		journeyRepository: journeyRepository,
		userRepository:    userRepository,
	}
}

// CreateOrder creates a new order or returns existing unpaid order
// Returns (result, isNewOrder, error) where isNewOrder is true for newly created orders
func (s *OrderService) CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*OrderResult, bool, error) {
	// 1. Validate that we have exactly one item (MVP constraint)
	if len(req.Items) != 1 {
		return nil, false, apperror.ValidationFailed()
	}

	item := req.Items[0]
	journeyID := item.JourneyID

	// 2. Check if user has already purchased this journey
	hasPurchased, err := s.orderRepository.CheckUserHasPurchasedJourney(ctx, userID, journeyID)
	if err != nil {
		return nil, false, apperror.DatabaseError(err)
	}
	if hasPurchased {
		return nil, false, apperror.JourneyAlreadyPurchased()
	}

	// 3. Check if user has existing unpaid order for this journey
	existingOrder, err := s.orderRepository.GetUnpaidOrderByUserAndJourney(ctx, userID, journeyID)
	if err != nil && !errors.Is(err, repository.ErrNoUnpaidOrderFound) {
		return nil, false, apperror.DatabaseError(err)
	}
	if err == nil {
		// Return existing order (not a new order)
		result, err := s.buildOrderResult(ctx, existingOrder)
		return result, false, err
	}

	// 4. Get journey to lock price
	journey, err := s.journeyRepository.GetByID(ctx, journeyID)
	if err != nil {
		if errors.Is(err, repository.ErrJourneyNotFound) {
			return nil, false, apperror.JourneyNotFound()
		}
		return nil, false, apperror.DatabaseError(err)
	}

	// 5. Calculate prices
	originalPrice := journey.Price
	discount := 0.0
	price := originalPrice - discount

	// 6. Generate order number: {timestamp(10)}{userId}{randomCode(5)}
	orderNumber, err := s.generateOrderNumber(userID)
	if err != nil {
		return nil, false, apperror.InternalServerError(err)
	}

	// 7. Create order
	expiredAt := time.Now().Add(72 * time.Hour) // 3 days
	order, err := s.orderRepository.CreateOrder(ctx, orderNumber, userID, originalPrice, discount, price, expiredAt)
	if err != nil {
		return nil, false, apperror.DatabaseError(err)
	}

	// 8. Create order item
	orderItem, err := s.orderRepository.CreateOrderItem(ctx, order.ID, journeyID, item.Quantity, originalPrice, discount, price)
	if err != nil {
		return nil, false, apperror.DatabaseError(err)
	}

	// 9. Build result
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, false, apperror.DatabaseError(err)
	}

	return &OrderResult{
		Order:         order,
		Items:         []model.OrderItem{*orderItem},
		JourneyTitles: map[int64]string{journeyID: journey.Title},
		Username:      user.Username,
	}, true, nil
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

	return s.buildOrderResult(ctx, order)
}

// buildOrderResult builds a complete order result with items and journey titles
func (s *OrderService) buildOrderResult(ctx context.Context, order *model.Order) (*OrderResult, error) {
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
