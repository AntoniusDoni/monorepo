package routes

import (
	core "github.com/antoniusDoni/monorepo/core/auth"
	"github.com/antoniusDoni/monorepo/shared/handler"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Register registers all shared routes and middleware
func Register(e *echo.Echo, authService *service.AuthService, db *gorm.DB, jwtSecret string, authMode string) {
	authHandler := handler.NewAuthHandler(authService)

	// Public routes
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected routes group
	apiGroup := e.Group("/api")

	// Initialize AuthMiddleware from core with JWT secret and DB
	authMiddleware := core.NewAuthMiddleware(jwtSecret, authMode, db)
	apiGroup.Use(authMiddleware.Middleware)

	// Example protected route with role middleware

}
