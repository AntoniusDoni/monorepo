package handler

import (
	"net/http"

	"github.com/antoniusDoni/monorepo/core/db/seeder"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// RunSeeder godoc
// @Summary      Run database seeder
// @Description  Executes the database seeder to populate initial data including sample warehouses, products, and unit products
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  object
// @Failure      401  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /admin/seed [post]
func (h *AdminHandler) RunSeeder(c echo.Context) error {
	// Run the seeder
	seeder.Seed(h.db)

	return c.JSON(http.StatusOK, contract.APIResponse[map[string]string]{
		Success: true,
		Data: map[string]string{
			"message": "Database seeder executed successfully",
			"status":  "completed",
		},
	})
}

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Returns the health status of the application
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  object
// @Router       /admin/health [get]
func (h *AdminHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, contract.APIResponse[map[string]string]{
		Success: true,
		Data: map[string]string{
			"status":  "healthy",
			"message": "Application is running properly",
		},
	})
}
