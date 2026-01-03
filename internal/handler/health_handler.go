package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool           *pgxpool.Pool
	logger         *slog.Logger
	requestTimeout time.Duration
}

func NewHealthHandler(pool *pgxpool.Pool, logger *slog.Logger, requestTimeout time.Duration) *HealthHandler {
	return &HealthHandler{
		pool:           pool,
		logger:         logger,
		requestTimeout: requestTimeout,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Create context with timeout for database ping
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	// Ping database to check connectivity
	if err := h.pool.Ping(ctx); err != nil {
		h.logger.Error("Health check failed: database ping error", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "DOWN",
			"database": "DOWN",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "UP",
		"database": "UP",
	})
}
