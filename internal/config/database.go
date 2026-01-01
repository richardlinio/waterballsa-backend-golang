package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func loadConfig() (*DatabaseConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Read from environment variables first, fallback to config.yaml
	host := getEnv("DB_HOST", viper.GetString("database.host"))
	portStr := getEnv("DB_PORT", strconv.Itoa(viper.GetInt("database.port")))
	user := getEnv("DB_USER", viper.GetString("database.user"))
	password := getEnv("DB_PASSWORD", viper.GetString("database.password"))
	name := getEnv("DB_NAME", viper.GetString("database.name"))
	sslmode := getEnv("DB_SSLMODE", viper.GetString("database.sslmode"))

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432 // default PostgreSQL port
	}

	config := &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslmode,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB() (*gorm.DB, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")

	return db, nil
}
