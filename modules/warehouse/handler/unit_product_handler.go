package handler

import (
	"net/http"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UnitProductHandler struct {
	service service.UnitProductService
}

func NewUnitProductHandler(service service.UnitProductService) *UnitProductHandler {
	return &UnitProductHandler{service: service}
}

func (h *UnitProductHandler) RegisterRoutes(g *echo.Group) {
	upg := g.Group("/unit-products")
	upg.GET("", h.GetAll)
	upg.POST("", h.Create)
	upg.GET("/:id", h.GetByID)
	upg.PUT("/:id", h.Update)
	upg.DELETE("/:id", h.Delete)
}

// GetAll godoc
// @Summary      Get list of unit products
// @Description  Retrieves paginated unit products optionally filtered by search term
// @Tags         unit-products
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        pageSize   query     int     false  "Page size (default: 10)"
// @Param        searchTerm query     string  false  "Search term to filter unit products by name"
// @Success      200       {object}  object
// @Failure      400       {object}  object
// @Failure      401       {object}  object
// @Failure      500       {object}  object
// @Security     BearerAuth
// @Router       /v1/api/unit-products [get]
func (h *UnitProductHandler) GetAll(c echo.Context) error {
	page := 1
	pageSize := 10
	searchTerm := c.QueryParam("searchTerm")

	if p := c.QueryParam("page"); p != "" {
		if parsedPage, err := parsePositiveInt(p); err == nil {
			page = parsedPage
		}
	}
	if ps := c.QueryParam("pageSize"); ps != "" {
		if parsedPageSize, err := parsePositiveInt(ps); err == nil {
			pageSize = parsedPageSize
		}
	}

	unitProducts, total, err := h.service.GetAll(page, pageSize, searchTerm)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := contract.PaginatedResponse[model.UnitProduct]{
		Items:      unitProducts,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return c.JSON(http.StatusOK, contract.APIResponse[contract.PaginatedResponse[model.UnitProduct]]{
		Success: true,
		Data:    resp,
	})
}

// Create godoc
// @Summary      Create a new unit product
// @Description  Create a new unit product with the provided information
// @Tags         unit-products
// @Accept       json
// @Produce      json
// @Param        unitProduct  body      model.UnitProduct  true  "Unit Product data"
// @Success      201       {object}  model.UnitProduct
// @Failure      400 {object} object
// @Failure      401 {object} object
// @Failure      500 {object} object
// @Security     BearerAuth
// @Router       /v1/api/unit-products [post]
func (h *UnitProductHandler) Create(c echo.Context) error {
	var unitProduct model.UnitProduct
	if err := c.Bind(&unitProduct); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Create(&unitProduct); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, contract.APIResponse[model.UnitProduct]{
		Success: true,
		Data:    unitProduct,
	})
}

// GetByID godoc
// @Summary      Get unit product by ID
// @Description  Retrieve a specific unit product by its ID
// @Tags         unit-products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Unit Product ID (UUID format)"
// @Success      200  {object}  model.UnitProduct
// @Failure      400 {object} object
// @Failure      401 {object} object
// @Failure      500 {object} object
// @Security     BearerAuth
// @Router       /v1/api/unit-products/{id} [get]
func (h *UnitProductHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	unitProduct, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}
	if unitProduct == nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   "unit product not found",
		})
	}
	return c.JSON(http.StatusOK, contract.APIResponse[model.UnitProduct]{
		Success: true,
		Data:    *unitProduct,
	})
}

// Update godoc
// @Summary      Update a unit product
// @Description  Update an existing unit product with new information
// @Tags         unit-products
// @Accept       json
// @Produce      json
// @Param        id          path      string            true  "Unit Product ID (UUID format)"
// @Param        unitProduct body      model.UnitProduct true  "Updated unit product data"
// @Success      200         {object}  model.UnitProduct
// @Failure      400 {object} object
// @Failure      401 {object} object
// @Failure      500 {object} object
// @Security     BearerAuth
// @Router       /v1/api/unit-products/{id} [put]
func (h *UnitProductHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	var unitProduct model.UnitProduct
	if err := c.Bind(&unitProduct); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Update(id, &unitProduct); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.UnitProduct]{
		Success: true,
		Data:    unitProduct,
	})
}

// Delete godoc
// @Summary      Delete a unit product
// @Description  Delete a unit product by its ID
// @Tags         unit-products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Unit Product ID (UUID format)"
// @Success      204 {string} string "No Content"
// @Failure      400 {object} object
// @Failure      401 {object} object
// @Failure      500 {object} object
// @Security     BearerAuth
// @Router       /v1/api/unit-products/{id} [delete]
func (h *UnitProductHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	if err := h.service.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
