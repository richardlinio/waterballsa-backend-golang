package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/dto"
	"github.com/richardlinio/waterballsa-backend-golang/internal/model"
	"github.com/richardlinio/waterballsa-backend-golang/internal/service"
)

type userService interface {
	GetProfile(ctx context.Context, userID int64) (*model.User, error)
}

type UserHandler struct {
	userService    userService
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewUserHandler(
	userService *service.UserService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

// GetCurrentUser handles GET /users/me
// Returns the authenticated user's profile
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Extract authenticated user from JWT middleware
	authenticatedUser := h.getAuthenticatedUser(c)
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

// getAuthenticatedUser extracts authenticated user from JWT context
func (h *UserHandler) getAuthenticatedUser(c *gin.Context) *model.User {
	payload, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return nil
	}

	user, ok := payload.(*model.User)
	if !ok {
		return nil
	}

	return user
}
