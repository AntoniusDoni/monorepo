package handler

import (
	"net/http"
	"strconv"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/labstack/echo/v4"
)

type WarehouseHandler struct {
	service service.WarehouseService
}

func NewWarehouseHandler(service service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) RegisterRoutes(g *echo.Group) {
	wg := g.Group("/warehouses")
	wg.GET("", h.GetAll)
	wg.POST("", h.Create)
	wg.GET("/:id", h.GetByID)
	wg.PUT("/:id", h.Update)
	wg.DELETE("/:id", h.Delete)
}

// GetAll godoc
// @Summary      Get list of warehouses
// @Description  Retrieves paginated warehouses optionally filtered by search term
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        pageSize  query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  object
// @Failure      400       {object}  object
// @Failure      401       {object}  object
// @Failure      500       {object}  object
// @Security     BearerAuth
// @Router       /v1/api/warehouses [get]
func (h *WarehouseHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	searchTerm := c.QueryParam("searchTerm")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	warehouses, total, err := h.service.GetAll(page, pageSize, searchTerm)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := contract.PaginatedResponse[model.Warehouse]{
		Items:      warehouses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return c.JSON(http.StatusOK, contract.APIResponse[contract.PaginatedResponse[model.Warehouse]]{
		Success: true,
		Data:    resp,
	})
}

// Create godoc
// @Summary      Create a new warehouse
// @Description  Create a new warehouse with the provided information
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        warehouse  body      model.Warehouse  true  "Warehouse data"
// @Success      201        {object}  model.Warehouse
// @Failure      400        {object}  object
// @Failure      401        {object}  object
// @Failure      500        {object}  object
// @Security     BearerAuth
// @Router       /v1/api/warehouses [post]
func (h *WarehouseHandler) Create(c echo.Context) error {
	var warehouse model.Warehouse
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Create(&warehouse); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, contract.APIResponse[model.Warehouse]{
		Success: true,
		Data:    warehouse,
	})
}

// GetByID godoc
// @Summary      Get warehouse by ID
// @Description  Retrieve a specific warehouse by its ID
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Warehouse ID"
// @Success      200  {object}  model.Warehouse
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/warehouses/{id} [get]
func (h *WarehouseHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id",
		})
	}

	warehouse, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}
	if warehouse == nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   "warehouse not found",
		})
	}
	return c.JSON(http.StatusOK, contract.APIResponse[model.Warehouse]{
		Success: true,
		Data:    *warehouse,
	})
}

// Update godoc
// @Summary      Update a warehouse
// @Description  Update an existing warehouse with new information
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id         path      int              true  "Warehouse ID"
// @Param        warehouse  body      model.Warehouse  true  "Updated warehouse data"
// @Success      200        {object}  model.Warehouse
// @Failure      400        {object}  object
// @Failure      401        {object}  object
// @Failure      500        {object}  object
// @Security     BearerAuth
// @Router       /v1/api/warehouses/{id} [put]
func (h *WarehouseHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id",
		})
	}

	var warehouse model.Warehouse
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := h.service.Update(uint(id), &warehouse); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.Warehouse]{
		Success: true,
		Data:    warehouse,
	})
}

// Delete godoc
// @Summary      Delete a warehouse
// @Description  Delete a warehouse by its ID
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Warehouse ID"
// @Success      204 {object}  object
// @Failure      400 {object}  object
// @Failure      401 {object}  object
// @Failure      500 {object}  object
// @Security     BearerAuth
// @Router       /v1/api/warehouses/{id} [delete]
func (h *WarehouseHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id",
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}