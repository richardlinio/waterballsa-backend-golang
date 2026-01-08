package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/handler"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/database"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/logger"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/server"
	"github.com/linporu/waterballsa-backend-golang/internal/repository"
	"github.com/linporu/waterballsa-backend-golang/internal/router"
	"github.com/linporu/waterballsa-backend-golang/internal/service"
	"github.com/linporu/waterballsa-backend-golang/internal/validator"
	"golang.org/x/sync/errgroup"
)

// ErrShutdownSignal indicates the application is shutting down due to OS signal
var ErrShutdownSignal = errors.New("received shutdown signal")

// Application manages the application lifecycle
type Application struct {
	config *config.Config
	db     *pgxpool.Pool
	server *server.HTTPServer
	logger *slog.Logger
}

// New creates and initializes a new Application
func New() (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Logger)

	// Register custom validators (must be done before any request handling starts)
	// These validators extend Gin's default validation rules for auth-specific fields
	// and are registered globally for the application lifetime
	if err := validator.RegisterAuthValidators(); err != nil {
		return nil, fmt.Errorf("failed to register validators: %w", err)
	}

	// Initialize database
	pool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize data layer (sqlc queries)
	queries := db.New(pool)

	// Initialize repository layer (data access)
	userRepository := repository.NewUserRepository(queries)

	// Initialize service layer (business logic)
	authService := service.NewAuthService(userRepository)

	// Initialize handler layer (HTTP handlers)
	healthHandler := handler.NewHealthHandler(pool, log, cfg.Server.RequestTimeout)
	authHandler := handler.NewAuthHandler(authService, log, cfg.Server.RequestTimeout)

	// Setup Gin router
	ginRouter := gin.Default()
	r := router.NewRouter(ginRouter, cfg.CORS, log, healthHandler, authHandler)
	r.Setup()

	// Create HTTP server
	httpServer := server.NewHTTPServer(cfg.Server, ginRouter)

	return &Application{
		config: cfg,
		db:     pool,
		server: httpServer,
		logger: log,
	}, nil
}

// Run starts the application and blocks until shutdown
func (a *Application) Run() error {
	a.logger.Info("Starting application...")

	// Create errgroup with context for coordinated shutdown
	g, ctx := errgroup.WithContext(context.Background())

	// Start HTTP server in errgroup
	g.Go(func() error {
		if err := a.server.Start(); err != nil {
			return fmt.Errorf("server failed: %w", err)
		}
		return nil
	})

	// Handle OS signals for graceful shutdown
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-quit:
			return fmt.Errorf("%w: %v", ErrShutdownSignal, sig)
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Wait for any goroutine to return
	err := g.Wait()

	// Log shutdown reason with context
	switch {
	case err == nil:
		a.logger.Info("Initiating graceful shutdown...")
	case errors.Is(err, ErrShutdownSignal):
		a.logger.Info("Received shutdown signal, initiating graceful shutdown...")
	default:
		a.logger.Error("Application error occurred, initiating shutdown...", "error", err)
	}

	// Perform graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if shutdownErr := a.server.Shutdown(shutdownCtx); shutdownErr != nil {
		a.logger.Error("Error shutting down HTTP server", "error", shutdownErr)
	}

	// Close database connection pool
	if a.db != nil {
		a.db.Close()
	}

	a.logger.Info("Application shutdown complete")

	// Return the original error that triggered shutdown (if it's not a signal)
	if err != nil && !errors.Is(err, ErrShutdownSignal) {
		return err
	}

	return nil
}
