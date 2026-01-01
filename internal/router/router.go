package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	healthHandler := handler.NewHealthHandler(db)

	r.GET("/healthz", healthHandler.HealthCheck)
}
