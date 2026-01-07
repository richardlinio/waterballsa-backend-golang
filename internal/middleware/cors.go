package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

// CORS creates and returns a CORS middleware with the given configuration
func CORS(config config.CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"Accept",
			"Origin",
			"X-Requested-With",
		},
		AllowCredentials: config.AllowCredentials,
		AllowWildcard:    true,
		MaxAge:           config.MaxAge,
	}
	return cors.New(corsConfig)
}
