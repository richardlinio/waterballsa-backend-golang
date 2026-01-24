package testutil

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/auth"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/database"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/logger"
	"github.com/linporu/waterballsa-backend-golang/internal/middleware"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"github.com/linporu/waterballsa-backend-golang/internal/router"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
	"github.com/linporu/waterballsa-backend-golang/internal/validator"
)

// TestServer encapsulates test server resources
type TestServer struct {
	Engine *gin.Engine
	Config *config.Config
	Pool   *pgxpool.Pool
	Logger *slog.Logger
	Server *http.Server
}

// NewTestServer creates and initializes a test server with all dependencies
// It uses the provided database connection information from Testcontainer
func NewTestServer(ctx context.Context, dbHost, dbPort string) (*TestServer, error) {
	// Set all required environment variables
	envVars := map[string]string{
		// Database configuration (from Testcontainer)
		"DB_HOST":     dbHost,
		"DB_PORT":     dbPort,
		"DB_USER":     TestDBUser,
		"DB_PASSWORD": TestDBPassword,
		"DB_NAME":     TestDBName,
		"DB_SSLMODE":  "disable",
		// Server configuration
		"SERVER_PORT":          "8081",
		"SERVER_HOST":          "127.0.0.1",
		"GIN_MODE":             "test",
		"LOG_LEVEL":            "info",
		"LOG_FORMAT":           "text",
		"CORS_ALLOWED_ORIGINS": "http://localhost:3000",
		"RATE_LIMIT_ENABLED":   "false",
		// JWT configuration
		"JWT_SECRET": "test-jwt-secret-key-at-least-32-characters-long-for-testing",
	}
	if err := setEnvVars(envVars); err != nil {
		return nil, err
	}

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Logger)

	// Register custom validators (idempotent - safe to call multiple times)
	if err := validator.RegisterAuthValidators(); err != nil {
		return nil, fmt.Errorf("failed to register validators: %w", err)
	}

	// Initialize database connection pool
	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Verify database connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize data layer (sqlc queries)
	queries := db.New(pool)

	// Initialize repository layer
	userRepository := repository.NewUserRepository(queries)
	accessTokenRepository := repository.NewAccessTokenRepository(queries)
	refreshTokenRepository := repository.NewRefreshTokenRepository(queries)
	journeyRepository := repository.NewJourneyRepository(queries)
	missionRepository := repository.NewMissionRepository(queries)
	progressRepository := repository.NewProgressRepository(queries)

	// Initialize token generator
	tokenGenerator := auth.NewJWTTokenGenerator(cfg.JWT)

	// Initialize service layer
	authService := service.NewAuthService(userRepository, accessTokenRepository, refreshTokenRepository, tokenGenerator)
	journeyService := service.NewJourneyService(journeyRepository)
	missionService := service.NewMissionService(missionRepository)
	progressService := service.NewProgressService(progressRepository, missionRepository)

	// Initialize JWT middleware
	jwtMiddleware, err := auth.NewJWTMiddleware(cfg.JWT, middleware.ExtractIdentity, middleware.Authorize, middleware.HandleUnauthorized)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWT middleware: %w", err)
	}

	// Initialize middleware (must be called to enable refresh token store)
	if err := jwtMiddleware.MiddlewareInit(); err != nil {
		return nil, fmt.Errorf("failed to initialize JWT middleware: %w", err)
	}

	// Initialize handler layer
	healthHandler := handler.NewHealthHandler(pool, log, cfg.Server.RequestTimeout)
	authHandler := handler.NewAuthHandler(authService, jwtMiddleware, cfg.JWT, log, cfg.Server.RequestTimeout)
	journeyHandler := handler.NewJourneyHandler(journeyService, log, cfg.Server.RequestTimeout)
	missionHandler := handler.NewMissionHandler(missionService, log, cfg.Server.RequestTimeout)
	progressHandler := handler.NewProgressHandler(progressService, log, cfg.Server.RequestTimeout)

	// Initialize blacklist checker middleware
	//nolint:contextcheck // False positive: middleware correctly captures context from c.Request.Context() at request time
	blacklistChecker := middleware.BlacklistChecker(accessTokenRepository)

	// Setup Gin router with test mode
	gin.SetMode(gin.TestMode)
	ginEngine := gin.New()
	r := router.NewRouter(ginEngine, cfg.CORS, cfg.RateLimit, cfg.JWT, log, healthHandler, authHandler, journeyHandler, missionHandler, progressHandler, jwtMiddleware, blacklistChecker)
	r.Setup()

	return &TestServer{
		Engine: ginEngine,
		Config: cfg,
		Pool:   pool,
		Logger: log,
	}, nil
}

// setEnvVars sets multiple environment variables and returns error on first failure
func setEnvVars(vars map[string]string) error {
	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}
	return nil
}

// Start starts the HTTP server in the background
// Returns the server instance for shutdown control
func (ts *TestServer) Start() error {
	addr := fmt.Sprintf("%s:%d", ts.Config.Server.Host, ts.Config.Server.Port)

	ts.Server = &http.Server{
		Addr:           addr,
		Handler:        ts.Engine,
		ReadTimeout:    ts.Config.Server.ReadTimeout,
		WriteTimeout:   ts.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in background goroutine
	go func() {
		if err := ts.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ts.Logger.Error("Test server failed to start", "error", err)
		}
	}()

	// Wait for server to be ready with simple retry
	healthURL := fmt.Sprintf("http://%s/healthz", addr)
	client := &http.Client{Timeout: 500 * time.Millisecond}

	for i := 0; i < 20; i++ { // Maximum wait 1 second (20 * 50ms)
		time.Sleep(50 * time.Millisecond)

		resp, err := client.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			if closeErr := resp.Body.Close(); closeErr != nil {
				ts.Logger.Warn("Failed to close health check response body", "error", closeErr)
			}
			ts.Logger.Info("Test server started", "address", addr)
			return nil
		}
		if resp != nil {
			if closeErr := resp.Body.Close(); closeErr != nil {
				ts.Logger.Warn("Failed to close health check response body", "error", closeErr)
			}
		}
	}

	return fmt.Errorf("test server failed to become ready within 1 second")
}

// Shutdown gracefully shuts down the test server
func (ts *TestServer) Shutdown(ctx context.Context) error {
	if ts.Server != nil {
		if err := ts.Server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown test server: %w", err)
		}
	}

	if ts.Pool != nil {
		ts.Pool.Close()
	}

	ts.Logger.Info("Test server shutdown complete")
	return nil
}

// BaseURL returns the base URL of the test server
func (ts *TestServer) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", ts.Config.Server.Host, ts.Config.Server.Port)
}
