package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
)

func SetupRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	healthHandler := handler.NewHealthHandler(pool)

	r.GET("/healthz", healthHandler.HealthCheck)
}
