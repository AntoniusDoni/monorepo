package config

import (
	"log"
	"os"
	"strconv"
	"time"
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
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	port, err := strconv.Atoi(getEnv("DB_Port", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_Port: %v", err)
	}

	expiredSec, err := strconv.Atoi(getEnv("Expired_Auth", "3000"))
	if err != nil {
		log.Fatalf("Invalid Expired_Auth: %v", err)
	}

	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBName:      getEnv("DB_Name", "monorepo"),
		DBUser:      getEnv("DB_User", "root"),
		DBPassword:  getEnv("DB_Password", "password"),
		DBPort:      port,
		DBDriver:    getEnv("DB_Driver", "postgresql"),
		DBTimeZone:  getEnv("DB_TimeZone", "Asia/Makassar"),
		JwtSecret:   getEnv("JWT_SECRET", ""),
		AuthExpired: time.Duration(expiredSec) * time.Second,
		AuthMode:    getEnv("AUTH_MODE", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
