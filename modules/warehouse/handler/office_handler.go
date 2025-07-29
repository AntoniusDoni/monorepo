package handler

import (
	"net/http"
	"strconv"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/labstack/echo/v4"
)

type OfficeHandler struct {
	service service.OfficeService
}

func NewOfficeHandler(service service.OfficeService) *OfficeHandler {
	return &OfficeHandler{service: service}
}
func (h *OfficeHandler) RegisterRoutes(g *echo.Group) {
	og := g.Group("/offices")
	og.GET("", h.GetAll)
	og.POST("", h.Create)
	og.GET("/:id", h.GetByID)
	og.PUT("/:id", h.Update)
	og.DELETE("/:id", h.Delete)
}
func (h *OfficeHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	searchTerm := c.QueryParam("searchTerm")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offices, total, err := h.service.GetAll(page, pageSize, searchTerm)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := contract.PaginatedResponse[model.Office]{
		Items:      offices,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return c.JSON(http.StatusOK, contract.APIResponse[contract.PaginatedResponse[model.Office]]{
		Success: true,
		Data:    resp,
	})
}
func (h *OfficeHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	office, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}
	if office == nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   "Office not found",
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.Office]{
		Success: true,
		Data:    *office,
	})
}
func (h *OfficeHandler) Create(c echo.Context) error {
	var office model.Office
	if err := c.Bind(&office); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Create(&office); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, contract.APIResponse[model.Office]{
		Success: true,
		Data:    office,
	})
}
func (h *OfficeHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var office model.Office
	if err := c.Bind(&office); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Update(id, &office); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.Office]{
		Success: true,
		Data:    office,
	})
}
func (h *OfficeHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, contract.APIResponse[any]{
		Success: true,
	})
}
