package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

// HTTPServer wraps the HTTP server with lifecycle management
type HTTPServer struct {
	server *http.Server
	config config.ServerConfig
}

// NewHTTPServer creates a new HTTP server with the given configuration and router
func NewHTTPServer(cfg config.ServerConfig, router *gin.Engine) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &HTTPServer{
		server: srv,
		config: cfg,
	}
}

// Start starts the HTTP server in a non-blocking way
func (s *HTTPServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	return nil
}
