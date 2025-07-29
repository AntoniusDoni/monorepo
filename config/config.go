package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBName      string
	DBUser      string
	DBPassword  string
	DBPort      int
	DBDriver    string
	DBTimeZone  string
	JwtSecret   string
	AuthExpired time.Duration
	AuthMode    string
	AppPort     string
}

func LoadConfig() (*Config, error) {
	// Load .env if available
	_ = godotenv.Load()

	dbPort, err := getEnvAsInt("DB_PORT", 5432)
	if err != nil {
		return nil, errors.New("invalid DB_PORT: " + err.Error())
	}

	authExpiredSec, err := getEnvAsInt("AUTH_EXPIRED", 3000)
	if err != nil {
		return nil, errors.New("invalid AUTH_EXPIRED: " + err.Error())
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET must be set")
	}

	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBName:      getEnv("DB_NAME", "monorepo"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBPort:      dbPort,
		DBDriver:    getEnv("DB_DRIVER", "postgresql"),
		DBTimeZone:  getEnv("DB_TIMEZONE", "Asia/Makassar"),
		JwtSecret:   jwtSecret,
		AuthExpired: time.Duration(authExpiredSec) * time.Second,
		AuthMode:    getEnv("AUTH_MODE", "jwt"),
		AppPort:     getEnv("APP_PORT", "8080"),
	}, nil
}

// getEnv returns the value of an environment variable or fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt converts an environment variable to int with fallback
func getEnvAsInt(key string, fallback int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback, nil
	}
	return strconv.Atoi(valueStr)
}
