package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
)

type Router struct {
	engine        *gin.Engine
	healthHandler *handler.HealthHandler
	authHandler   *handler.AuthHandler
}

func NewRouter(
	engine *gin.Engine,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) *Router {
	return &Router{
		engine:        engine,
		healthHandler: healthHandler,
		authHandler:   authHandler,
	}
}

func (r *Router) Setup() {
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
