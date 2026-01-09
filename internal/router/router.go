package router

import (
	"log/slog"

	"github.com/gin-contrib/requestid"
	ginslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"github.com/linporu/waterballsa-backend-golang/internal/middleware"
)

type Router struct {
	engine        *gin.Engine
	corsConfig    config.CORSConfig
	logger        *slog.Logger
	healthHandler *handler.HealthHandler
	authHandler   *handler.AuthHandler
}

func NewRouter(
	engine *gin.Engine,
	corsConfig config.CORSConfig,
	logger *slog.Logger,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) *Router {
	return &Router{
		engine:        engine,
		corsConfig:    corsConfig,
		logger:        logger,
		healthHandler: healthHandler,
		authHandler:   authHandler,
	}
}

func (r *Router) Setup() {
	// 1. Setup middlewares
	r.engine.Use(gin.Recovery())
	r.setupRequestIDMiddleware()
	r.setupLoggingMiddleware()
	r.engine.Use(middleware.Security())
	r.engine.Use(middleware.CORS(r.corsConfig))
	r.engine.Use(middleware.ErrorHandler(r.logger))

	// 2. Setup routes
	r.setupHealthRoutes()
	r.setupAuthRoutes()
}

func (r *Router) setupRequestIDMiddleware() {
	r.engine.Use(requestid.New(
		requestid.WithHandler(func(c *gin.Context, id string) {
			c.Header("X-Request-ID", id)
		})))
}

func (r *Router) setupLoggingMiddleware() {
	r.engine.Use(ginslog.SetLogger(
		ginslog.WithLogger(func(c *gin.Context, l *slog.Logger) *slog.Logger {
			return r.logger.With("request_id", requestid.Get(c))
		}),

		ginslog.WithSkipPath([]string{"/healthz"}),

		ginslog.WithSkipper(func(c *gin.Context) bool {
			// Skip request logging for errors (ErrorHandler logs them with more detail)
			return len(c.Errors) > 0
		}),
	))
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
