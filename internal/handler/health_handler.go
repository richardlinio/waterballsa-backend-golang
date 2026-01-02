package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Create context with timeout for database ping
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping database to check connectivity
	if err := h.pool.Ping(ctx); err != nil {
		log.Printf("Health check failed: database ping error: %v", err)
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
