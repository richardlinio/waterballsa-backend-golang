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

type userService interface {
	GetProfile(ctx context.Context, userID int64) (*model.User, error)
}

type userJourneyService interface {
	GetUserJourneys(ctx context.Context, userID, authenticatedUserID int64) ([]dto.UserJourneyItem, error)
}

type UserHandler struct {
	userService        userService
	userJourneyService userJourneyService
	logger             *slog.Logger
	requestTimeout     time.Duration
}

func NewUserHandler(
	userService *service.UserService,
	userJourneyService *service.UserJourneyService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *UserHandler {
	return &UserHandler{
		userService:        userService,
		userJourneyService: userJourneyService,
		logger:             logger,
		requestTimeout:     requestTimeout,
	}
}

// GetCurrentUser handles GET /users/me
// Returns the authenticated user's profile
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Extract authenticated user from JWT middleware
	authenticatedUser := util.GetAuthenticatedUser(c)
	if authenticatedUser == nil {
		_ = c.Error(apperror.Unauthorized())
		return
	}

	// Get user profile from service
	user, err := h.userService.GetProfile(ctx, authenticatedUser.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToUserProfileResponse(user)
	c.JSON(http.StatusOK, response)
}

// GetUserJourneys handles GET /users/:userId/journeys
func (h *UserHandler) GetUserJourneys(c *gin.Context) {
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Call service with authorization check
	journeys, err := h.userJourneyService.GetUserJourneys(ctx, userID, authenticatedUser.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToUserJourneyListResponse(journeys)
	c.JSON(http.StatusOK, response)
}
