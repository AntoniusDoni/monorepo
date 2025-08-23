package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/antoniusDoni/monorepo/config"
	"github.com/antoniusDoni/monorepo/core"
	_ "github.com/antoniusDoni/monorepo/docs"

	auth "github.com/antoniusDoni/monorepo/core/auth"
	dbpkg "github.com/antoniusDoni/monorepo/core/db"

	modules "github.com/antoniusDoni/monorepo/modules"
	wrepo "github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	shandler "github.com/antoniusDoni/monorepo/shared/handler"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/routes"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

const (
	// Default server configuration constants
	DefaultPort                   = "8080"
	DefaultRequestTimeoutSeconds  = 30
	DefaultShutdownTimeoutSeconds = 10
)

// @title           Your Monorepo API
// @version         1.0
// @description     API documentation for your monorepo services
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize configuration
	cfg, err := initializeConfig()
	if err != nil {
		log.Fatal("Failed to initialize configuration:", err)
	}

	// Initialize database
	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize shared services
	sharedServices, err := initializeSharedServices(dbInstance, cfg)
	if err != nil {
		log.Fatal("Failed to initialize shared services:", err)
	}

	// Initialize Echo server
	e := initializeEchoServer()

	// Register Swagger documentation FIRST (before any auth middleware)
	registerSwaggerRoutes(e, getServerPort(cfg))

	// Register shared routes
	routes.Register(e, sharedServices.AuthService, dbInstance, cfg.JwtSecret, cfg.AuthMode)

	// Register admin routes
	registerAdminRoutes(e, dbInstance, cfg, sharedServices.AuthService)

	// Initialize and register modules using the new registry system
	if err := initializeModules(e, sharedServices, cfg); err != nil {
		log.Fatal("Failed to initialize modules:", err)
	}

	// Start server with graceful shutdown
	startServerWithGracefulShutdown(e, getServerPort(cfg))
}

// initializeConfig loads and validates configuration
func initializeConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

// initializeDatabase establishes database connection
func initializeDatabase() (*gorm.DB, error) {
	dbInstance, err := dbpkg.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return dbInstance, nil
}

// initializeSharedServices creates shared service instances
func initializeSharedServices(db *gorm.DB, cfg *config.Config) (*modules.ModuleContext, error) {
	userRepo := repository.NewUserRepository(db)
	officeRepo := wrepo.NewOfficeRepository(db)
	authService := service.NewAuthService(userRepo, officeRepo, db, cfg.JwtSecret)

	modCtx := &modules.ModuleContext{
		DB:          db,
		UserRepo:    userRepo,
		AuthService: authService,
		OfficeRepo:  officeRepo,
	}

	log.Println("Shared services initialized successfully")
	return modCtx, nil
}

// initializeEchoServer creates and configures Echo server
func initializeEchoServer() *echo.Echo {
	e := echo.New()

	// Set custom validator
	e.Validator = core.NewValidator()

	// Add middleware with Swagger route skipper
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: swaggerSkipper,
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper: swaggerSkipper,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: swaggerSkipper,
	}))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Skipper: swaggerSkipper,
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: swaggerSkipper,
		Timeout: time.Duration(getRequestTimeout()) * time.Second,
	}))

	log.Printf("Echo server initialized with middleware (request timeout: %ds)", getRequestTimeout())
	return e
}

// swaggerSkipper skips middleware for Swagger routes
func swaggerSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	return strings.HasPrefix(path, "/swagger")
}

// initializeModules loads and registers enabled modules using the registry system
func initializeModules(e *echo.Echo, services *modules.ModuleContext, cfg *config.Config) error {
	// Create API group with authentication middleware
	apiGroup := e.Group("/v1/api")
	authMiddleware := auth.NewAuthMiddleware(cfg.JwtSecret, cfg.AuthMode, services.DB)
	apiGroup.Use(authMiddleware.Middleware)

	// Setup module registry and register all available modules
	registry := modules.SetupModules(services)

	// Initialize all enabled modules
	return registry.InitializeModules(apiGroup, services)
}

// registerAdminRoutes registers admin endpoints
func registerAdminRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, authService *service.AuthService) {
	// Create admin handler
	adminHandler := shandler.NewAdminHandler(db)

	// Create admin group with authentication middleware
	adminGroup := e.Group("/admin")
	authMiddleware := auth.NewAuthMiddleware(cfg.JwtSecret, cfg.AuthMode, db)
	adminGroup.Use(authMiddleware.Middleware)

	// Register admin routes
	adminGroup.POST("/seed", adminHandler.RunSeeder)

	// Health check endpoint (no auth required)
	e.GET("/admin/health", adminHandler.HealthCheck)

	log.Println("Admin routes registered successfully")
}

// getServerPort determines the server port from config or environment
func getServerPort(cfg *config.Config) string {
	// Check environment variable first (highest priority)
	if port := os.Getenv("PORT"); port != "" {
		if _, err := strconv.Atoi(port); err == nil {
			return port
		}
		log.Printf("Invalid PORT environment variable '%s', using default %s", port, DefaultPort)
	}

	// Default port
	return DefaultPort
}

// getRequestTimeout gets request timeout from environment or returns default
func getRequestTimeout() int {
	if timeout := os.Getenv("REQUEST_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil && val > 0 {
			return val
		}
		log.Printf("Invalid REQUEST_TIMEOUT environment variable '%s', using default %d", timeout, DefaultRequestTimeoutSeconds)
	}
	return DefaultRequestTimeoutSeconds
}

// getShutdownTimeout gets shutdown timeout from environment or returns default
func getShutdownTimeout() int {
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil && val > 0 {
			return val
		}
		log.Printf("Invalid SHUTDOWN_TIMEOUT environment variable '%s', using default %d", timeout, DefaultShutdownTimeoutSeconds)
	}
	return DefaultShutdownTimeoutSeconds
}

// registerSwaggerRoutes registers Swagger documentation routes with dynamic configuration
func registerSwaggerRoutes(e *echo.Echo, port string) {
	// Register Swagger routes directly on the main Echo instance (bypassing any group middleware)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add a redirect from /swagger to /swagger/
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/")
	})

	// Add a test endpoint to verify Swagger is working
	e.GET("/swagger/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Swagger routes are working without authentication",
			"status":  "success",
		})
	})

	// Add Swagger JSON endpoint explicitly
	e.GET("/swagger/doc.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Swagger JSON endpoint is accessible",
			"status":  "success",
		})
	})

	log.Printf("Swagger documentation registered at: http://localhost:%s/swagger/ (no auth required)", port)
}

// startServerWithGracefulShutdown starts the server and handles graceful shutdown
func startServerWithGracefulShutdown(e *echo.Echo, port string) {
	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		log.Printf("Invalid port '%s', using default %s", port, DefaultPort)
		port = DefaultPort
	}

	// Start server in a goroutine
	go func() {
		address := ":" + port
		log.Printf("Starting server on %s", address)
		log.Printf("Swagger documentation available at: http://localhost:%s/swagger/", port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownTimeout := getShutdownTimeout()
	log.Printf("Shutting down server (timeout: %ds)...", shutdownTimeout)

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
