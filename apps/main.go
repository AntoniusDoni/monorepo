package main

import (
	"log"

	"github.com/antoniusDoni/monorepo/config"
	"github.com/antoniusDoni/monorepo/core"
	dbpkg "github.com/antoniusDoni/monorepo/core/db"
	Seeder "github.com/antoniusDoni/monorepo/core/db/seeder"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/routes"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

	// Load config (e.g. DB config, JWT secret)
	cfg := config.LoadConfig()

	// Initialize DB instance
	dbInstance, err := dbpkg.GetInstance()
	if err != nil {
		log.Fatal("Failed to get DB instance:", err)
	}

	// Seed initial data (optional)
	Seeder.Seed(dbInstance)

	// Initialize user repository and auth service
	userRepo := repository.NewUserRepository(dbInstance)
	authService := service.NewAuthService(userRepo, cfg.JwtSecret)

	// Initialize Echo
	e := echo.New()
	e.Validator = core.NewValidator()

	// Register all shared routes, injecting auth service, DB, JWT secret and auth mode
	routes.Register(e, authService, dbInstance, cfg.JwtSecret, "jwt")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
