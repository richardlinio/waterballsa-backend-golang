package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func loadConfig() (*DatabaseConfig, error) {
	// Try to read config file, but don't fail if it doesn't exist
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Ignore error if config file doesn't exist
	_ = viper.ReadInConfig()

	// Read from environment variables first, fallback to config.yaml
	host := getEnv("DB_HOST", viper.GetString("database.host"))
	portStr := getEnv("DB_PORT", strconv.Itoa(viper.GetInt("database.port")))
	user := getEnv("DB_USER", viper.GetString("database.user"))
	password := getEnv("DB_PASSWORD", viper.GetString("database.password"))
	name := getEnv("DB_NAME", viper.GetString("database.name"))
	sslmode := getEnv("DB_SSLMODE", viper.GetString("database.sslmode"))

	maxConnsStr := getEnv("DB_MAX_CONNS", strconv.Itoa(viper.GetInt("database.max_conns")))
	minConnsStr := getEnv("DB_MIN_CONNS", strconv.Itoa(viper.GetInt("database.min_conns")))

	// Validate required fields
	if host == "" {
		return nil, fmt.Errorf("database host is required (set DB_HOST env var or database.host in config)")
	}
	if user == "" {
		return nil, fmt.Errorf("database user is required (set DB_USER env var or database.user in config)")
	}
	if name == "" {
		return nil, fmt.Errorf("database name is required (set DB_NAME env var or database.name in config)")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432 // default PostgreSQL port
	}

	// Parse pool settings with defaults
	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil || maxConns == 0 {
		maxConns = 10 // default
	}

	minConns, err := strconv.Atoi(minConnsStr)
	if err != nil || minConns == 0 {
		minConns = 2 // default
	}

	config := &DatabaseConfig{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		Name:            name,
		SSLMode:         sslmode,
		MaxConns:        int32(maxConns),
		MinConns:        int32(minConns),
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB() (*pgxpool.Pool, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Build connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)

	// Configure pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Create context with timeout for initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	fmt.Println("Database connection pool established successfully")

	return pool, nil
}
