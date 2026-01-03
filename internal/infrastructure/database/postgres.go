package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linporu/waterballsa-backend-golang/internal/config"
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Configure pool using Config struct to avoid exposing password in connection string
	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to create pool config: %w", err)
	}

	// Set connection parameters securely
	poolConfig.ConnConfig.Host = cfg.Host
	poolConfig.ConnConfig.Port = uint16(cfg.Port)
	poolConfig.ConnConfig.User = cfg.User
	poolConfig.ConnConfig.Password = cfg.Password
	poolConfig.ConnConfig.Database = cfg.Name

	// Configure TLS based on SSL mode
	// Note: sslmode is a client-side parameter, not a server runtime parameter
	if cfg.SSLMode == "disable" {
		poolConfig.ConnConfig.TLSConfig = nil
	}
	// For other SSL modes (require, verify-ca, verify-full), TLS config would be needed
	// Currently only 'disable' mode is fully supported for development environment

	// Set pool settings
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// Create context with timeout for initial connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Connect to database
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
