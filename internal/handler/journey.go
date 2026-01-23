package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

// journeyService defines the journey service operations needed by the handler
type journeyService interface {
	List(ctx context.Context) ([]*model.Journey, error)
}

type JourneyHandler struct {
	journeyService journeyService
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewJourneyHandler(
	journeyService *service.JourneyService,
	logger *slog.Logger,
	requestTimeout time.Duration,
) *JourneyHandler {
	return &JourneyHandler{
		journeyService: journeyService,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *JourneyHandler) GetJourneys(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	journeys, err := h.journeyService.List(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Convert domain models to DTOs
	items := make([]dto.JourneyListItem, 0, len(journeys))
	for _, journey := range journeys {
		items = append(items, dto.JourneyListItem{
			ID:            journey.ID,
			Slug:          journey.Slug,
			Title:         journey.Title,
			Description:   journey.Description,
			CoverImageURL: journey.CoverImageURL,
			TeacherName:   journey.TeacherName,
			Price:         journey.Price,
		})
	}

	c.JSON(http.StatusOK, dto.JourneyListResponse{
		Journeys: items,
	})
}
