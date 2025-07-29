package routes

import (
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
}
