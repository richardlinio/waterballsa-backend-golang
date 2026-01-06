package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
)

func SetupRoutes(
	r *gin.Engine,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) {
	r.GET("/healthz", healthHandler.HealthCheck)
	r.POST("/auth/register", authHandler.Register)
}
