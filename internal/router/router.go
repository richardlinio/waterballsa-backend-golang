package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
)

func SetupRoutes(
	r *gin.Engine,
	healthHandler *handler.HealthHandler,
	registerHandler *handler.RegisterHandler,
) {
	r.GET("/healthz", healthHandler.HealthCheck)
	r.POST("/auth/register", registerHandler.Register)
}
