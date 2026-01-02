package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linporu/waterballsa-backend-golang/internal/app"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	// Initialize application
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start application in goroutine
	go func() {
		if err := application.Run(); err != nil {
			log.Fatal("Application run failed:", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown failed:", err)
	}

	log.Println("Application exited")
}
