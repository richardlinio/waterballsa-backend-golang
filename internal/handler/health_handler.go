package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(503, gin.H{
			"status":   "DOWN",
			"database": "DOWN",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(503, gin.H{
			"status":   "DOWN",
			"database": "DOWN",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "UP",
		"database": "UP",
	})
}
