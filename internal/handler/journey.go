package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/apperror"
	"github.com/linporu/waterballsa-backend-golang/internal/dto"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
)

// journeyService defines the journey service operations needed by the handler
type journeyService interface {
	List(ctx context.Context) ([]*model.Journey, error)
	GetDetail(ctx context.Context, journeyID int64) (*model.Journey, []*model.Chapter, map[int64][]*model.Mission, error)
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

	response := dto.ToJourneyListResponse(journeys)
	c.JSON(http.StatusOK, response)
}

func (h *JourneyHandler) GetJourneyDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Parse journey ID from path parameter
	journeyIDStr := c.Param("journeyId")
	journeyID, err := strconv.ParseInt(journeyIDStr, 10, 64)
	if err != nil {
		_ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
		return
	}

	// Get journey detail from service
	journey, chapters, missionsByChapter, err := h.journeyService.GetDetail(ctx, journeyID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := dto.ToJourneyDetailResponse(journey, chapters, missionsByChapter)
	c.JSON(http.StatusOK, response)
}
