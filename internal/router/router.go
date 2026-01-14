package router

import (
	"log/slog"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-contrib/requestid"
	ginslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"github.com/linporu/waterballsa-backend-golang/internal/middleware"
)

type Router struct {
	engine          *gin.Engine
	corsConfig      config.CORSConfig
	rateLimitConfig config.RateLimitConfig
	logger          *slog.Logger
	healthHandler   *handler.HealthHandler
	authHandler     *handler.AuthHandler
	jwtMiddleware   *jwt.GinJWTMiddleware
}

func NewRouter(
	engine *gin.Engine,
	corsConfig config.CORSConfig,
	rateLimitConfig config.RateLimitConfig,
	logger *slog.Logger,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	jwtMiddleware *jwt.GinJWTMiddleware,
) *Router {
	return &Router{
		engine:          engine,
		corsConfig:      corsConfig,
		rateLimitConfig: rateLimitConfig,
		logger:          logger,
		healthHandler:   healthHandler,
		authHandler:     authHandler,
		jwtMiddleware:   jwtMiddleware,
	}
}

func (r *Router) Setup() {
	// 1. Setup middlewares
	r.engine.Use(gin.Recovery())
	r.setupRequestIDMiddleware()
	r.setupLoggingMiddleware()
	r.engine.Use(middleware.ErrorHandler(r.logger))
	r.engine.Use(middleware.CORS(r.corsConfig))
	r.engine.Use(middleware.RateLimit(r.rateLimitConfig, r.logger))
	r.engine.Use(middleware.Security())

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
		// Public routes
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.jwtMiddleware.LoginHandler)
		auth.POST("/refresh", r.jwtMiddleware.RefreshHandler)
	}

	// Protected routes (require JWT authentication)
	authProtected := r.engine.Group("/auth")
	authProtected.Use(r.jwtMiddleware.MiddlewareFunc())
	{
		authProtected.POST("/logout", r.jwtMiddleware.LogoutHandler)
	}
}
