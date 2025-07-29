package main

import (
	"log"
	"os"
	"strings"

	"github.com/antoniusDoni/monorepo/config"
	"github.com/antoniusDoni/monorepo/core"
	auth "github.com/antoniusDoni/monorepo/core/auth"
	dbpkg "github.com/antoniusDoni/monorepo/core/db"
	Seeder "github.com/antoniusDoni/monorepo/core/db/seeder"
	modules "github.com/antoniusDoni/monorepo/modules"
	"github.com/antoniusDoni/monorepo/modules/warehouse/handler"
	wrepo "github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	wservice "github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/routes"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dbInstance, err := dbpkg.GetInstance()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	Seeder.Seed(dbInstance)

	// Initialize shared repositories/services
	userRepo := repository.NewUserRepository(dbInstance)
	authService := service.NewAuthService(userRepo, cfg.JwtSecret)

	modCtx := modules.ModuleContext{
		DB:          dbInstance,
		UserRepo:    userRepo,
		AuthService: authService,
	}

	e := echo.New()
	e.Validator = core.NewValidator()
	// Register shared routes like login, register, and middleware
	routes.Register(e, authService, dbInstance, cfg.JwtSecret, cfg.AuthMode)

	enabledModules := []handler.RouteRegistrar{}

	enabled := map[string]bool{}
	for _, m := range strings.Split(os.Getenv("ENABLE_MODULES"), ",") {
		enabled[strings.TrimSpace(m)] = true
	}

	if enabled["warehouse"] {
		// Create warehouse repo with shared DB
		whRepo := wrepo.NewWarehouseRepository(modCtx.DB)

		// Create warehouse service with warehouse repo
		whService := wservice.NewWarehouseService(whRepo)

		// Create handler with service
		whHandler := handler.NewWarehouseHandler(whService)

		enabledModules = append(enabledModules, whHandler)
	}

	// Register all enabled module routes
	apiGroup := e.Group("/v1/api")

	// Initialize AuthMiddleware from core with JWT secret and DB
	authMiddleware := auth.NewAuthMiddleware(cfg.JwtSecret, cfg.AuthMode, dbInstance)
	apiGroup.Use(authMiddleware.Middleware)
	for _, mod := range enabledModules {
		mod.RegisterRoutes(apiGroup)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
