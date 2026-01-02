package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/database"
	"github.com/linporu/waterballsa-backend-golang/internal/infrastructure/server"
	"github.com/linporu/waterballsa-backend-golang/internal/router"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout = 5 * time.Second
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

// Run starts the application and blocks until shutdown
func (a *Application) Run() error {
	log.Println("Starting application...")

	// Create errgroup with context for coordinated shutdown
	g, ctx := errgroup.WithContext(context.Background())

	// Start HTTP server in errgroup
	g.Go(func() error {
		log.Println("HTTP server goroutine started")
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
			log.Printf("Received signal: %v", sig)
			return fmt.Errorf("received shutdown signal: %v", sig)
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Wait for any goroutine to return
	err := g.Wait()

	// Perform graceful shutdown
	log.Println("Shutting down application...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if shutdownErr := a.server.Shutdown(shutdownCtx); shutdownErr != nil {
		log.Printf("Error shutting down HTTP server: %v", shutdownErr)
	}

	// Close database connection pool
	if a.db != nil {
		a.db.Close()
		log.Println("Database connection pool closed")
	}

	log.Println("Application shutdown complete")

	// Return the original error that triggered shutdown (if it's not a signal)
	if err != nil && err.Error() != "received shutdown signal: interrupt" && err.Error() != "received shutdown signal: terminated" {
		return err
	}

	return nil
}
