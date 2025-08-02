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
	og.GET("/active", h.GetActiveOffices)
	og.POST("", h.Create)
	og.GET("/:id", h.GetByID)
	og.PUT("/:id", h.Update)
	og.DELETE("/:id", h.Delete)
}

// GetAll godoc
// @Summary      Get list of offices
// @Description  Retrieves paginated offices optionally filtered by search term
// @Tags         offices
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        pageSize   query     int     false  "Page size"
// @Param        searchTerm query     string  false  "Search term"
// @Success      200        {object}  object
// @Failure      401        {object}  object
// @Failure      500        {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices [get]
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

// GetByID godoc
// @Summary      Get office by ID
// @Description  Retrieve a specific office by its ID
// @Tags         offices
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Office ID"
// @Success      200  {object}  model.Office
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices/{id} [get]
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

// Create godoc
// @Summary      Create a new office
// @Description  Create a new office with the provided information
// @Tags         offices
// @Accept       json
// @Produce      json
// @Param        office  body      model.Office  true  "Office data"
// @Success      201     {object}  model.Office
// @Failure      400     {object}  object
// @Failure      401     {object}  object
// @Failure      500     {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices [post]
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

// Update godoc
// @Summary      Update an office
// @Description  Update an existing office with new information
// @Tags         offices
// @Accept       json
// @Produce      json
// @Param        id      path      string        true  "Office ID"
// @Param        office  body      model.Office  true  "Updated office data"
// @Success      200     {object}  model.Office
// @Failure      400     {object}  object
// @Failure      401     {object}  object
// @Failure      500     {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices/{id} [put]
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

// Delete godoc
// @Summary      Delete an office
// @Description  Delete an office by its ID
// @Tags         offices
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Office ID"
// @Success      204 {object}  object
// @Failure      401 {object}  object
// @Failure      500 {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices/{id} [delete]
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

// GetActiveOffices godoc
// @Summary      Get active offices
// @Description  Retrieve all offices with active status
// @Tags         offices
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.Office
// @Failure      401  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/offices/active [get]
func (h *OfficeHandler) GetActiveOffices(c echo.Context) error {
	offices, err := h.service.GetActiveOffices()
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, contract.APIResponse[[]model.Office]{
		Success: true,
		Data:    offices,
	})
}