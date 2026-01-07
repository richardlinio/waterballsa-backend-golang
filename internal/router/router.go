package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"github.com/linporu/waterballsa-backend-golang/internal/middleware"
)

type Router struct {
	engine        *gin.Engine
	corsConfig    config.CORSConfig
	healthHandler *handler.HealthHandler
	authHandler   *handler.AuthHandler
}

func NewRouter(
	engine *gin.Engine,
	corsConfig config.CORSConfig,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) *Router {
	return &Router{
		engine:        engine,
		corsConfig:    corsConfig,
		healthHandler: healthHandler,
		authHandler:   authHandler,
	}
}

func (r *Router) Setup() {
	// 1. Setup middlewares
	r.engine.Use(middleware.CORS(r.corsConfig))

	// 2. Setup routes
	r.setupHealthRoutes()
	r.setupAuthRoutes()
}

func (r *Router) setupHealthRoutes() {
	r.engine.GET("/healthz", r.healthHandler.HealthCheck)
}

func (r *Router) setupAuthRoutes() {
	auth := r.engine.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
	}
}
