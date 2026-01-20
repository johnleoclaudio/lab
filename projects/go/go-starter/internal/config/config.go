package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server Configuration
	ServerAddress string
	ServerEnv     string

	// Database Configuration
	DatabaseURL                   string
	DatabaseMaxConnections        int
	DatabaseMaxIdleConnections    int
	DatabaseConnectionMaxLifetime time.Duration

	// JWT Configuration
	JWTSecret        string
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration

	// Redis Configuration
	RedisURL string

	// Logging Configuration
	LogLevel  string
	LogFormat string

	// CORS Configuration
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		ServerEnv:     getEnv("SERVER_ENV", "development"),

		DatabaseURL:                   getEnv("DATABASE_URL", ""),
		DatabaseMaxConnections:        getEnvInt("DATABASE_MAX_CONNECTIONS", 25),
		DatabaseMaxIdleConnections:    getEnvInt("DATABASE_MAX_IDLE_CONNECTIONS", 10),
		DatabaseConnectionMaxLifetime: getEnvDuration("DATABASE_CONNECTION_MAX_LIFETIME", 5*time.Minute),

		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTExpiry:        getEnvDuration("JWT_EXPIRY", 24*time.Hour),
		JWTRefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),

		RedisURL: getEnv("REDIS_URL", ""),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
