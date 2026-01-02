package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
)

func SetupRoutes(r *gin.Engine, pool *pgxpool.Pool, logger *slog.Logger) {
	healthHandler := handler.NewHealthHandler(pool, logger)

	r.GET("/healthz", healthHandler.HealthCheck)
}
