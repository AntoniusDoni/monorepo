package warehouse

import (
	"log"

	"github.com/antoniusDoni/monorepo/modules/warehouse/handler"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ModuleDependencies holds the dependencies needed for the warehouse module
type ModuleDependencies struct {
	DB         *gorm.DB
	OfficeRepo repository.OfficeRepository
}

// RegisterRoutes registers all warehouse module routes
func RegisterRoutes(apiGroup *echo.Group, deps *ModuleDependencies) error {
	// Initialize warehouse handler
	whRepo := repository.NewWarehouseRepository(deps.DB)
	whService := service.NewWarehouseService(whRepo)
	whHandler := handler.NewWarehouseHandler(whService)

	// Initialize office handler
	officeService := service.NewOfficeService(deps.OfficeRepo)
	officeHandler := handler.NewOfficeHandler(officeService)

	// Initialize category product repository (shared)
	categoryProductRepo := repository.NewCategoryProductRepository(deps.DB)

	// Initialize product handler
	productRepo := repository.NewProductRepository(deps.DB)
	productService := service.NewProductService(productRepo, categoryProductRepo)
	productHandler := handler.NewProductHandler(productService)

	// Initialize unit product handler
	unitProductRepo := repository.NewUnitProductRepository(deps.DB)
	unitProductService := service.NewUnitProductService(unitProductRepo)
	unitProductHandler := handler.NewUnitProductHandler(unitProductService)

	// Initialize category product handler
	categoryProductService := service.NewCategoryProductService(categoryProductRepo)
	categoryProductHandler := handler.NewCategoryProductHandler(categoryProductService)

	// Register all handlers
	handlers := []handler.RouteRegistrar{
		whHandler,
		officeHandler,
		productHandler,
		unitProductHandler,
		categoryProductHandler,
	}

	for _, h := range handlers {
		h.RegisterRoutes(apiGroup)
	}

	log.Printf("Warehouse module: registered %d handlers", len(handlers))
	return nil
}