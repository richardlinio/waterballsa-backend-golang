package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/database"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/server"
	"github.com/linporu/waterballsa-backend-golang/internal/router"
)

// Application manages the application lifecycle
type Application struct {
	config *config.Config
	db     *pgxpool.Pool
	server *server.HTTPServer
}

// New creates and initializes a new Application
func New() (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database
	pool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Setup Gin router
	ginRouter := gin.Default()
	router.SetupRoutes(ginRouter, pool)

	// Create HTTP server
	httpServer := server.NewHTTPServer(cfg.Server, ginRouter)

	return &Application{
		config: cfg,
		db:     pool,
		server: httpServer,
	}, nil
}

// Run starts the application
func (a *Application) Run() error {
	log.Println("Starting application...")

	// Start HTTP server (blocking)
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown(ctx context.Context) error {
	log.Println("Shutting down application...")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down HTTP server: %v", err)
	}

	// Close database connection pool
	if a.db != nil {
		a.db.Close()
		log.Println("Database connection pool closed")
	}

	log.Println("Application shutdown complete")
	return nil
}
