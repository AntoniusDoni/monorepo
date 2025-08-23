package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/antoniusDoni/monorepo/modules/warehouse/dto"
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(g *echo.Group) {
	pg := g.Group("/products")
	pg.GET("", h.GetAll)
	pg.POST("", h.Create)
	pg.GET("/:id", h.GetByID)
	pg.PUT("/:id", h.Update)
	pg.DELETE("/:id", h.Delete)
}

// GetAll godoc
// @Summary      Get list of products
// @Description  Retrieves paginated products optionally filtered by search term
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        pageSize   query     int     false  "Page size (default: 10)"
// @Param        searchTerm query     string  false  "Search term to filter products by name or description"
// @Success      200        {object}  object
// @Failure      400        {object}  object
// @Failure      401        {object}  object
// @Failure      500        {object}  object
// @Security     BearerAuth
// @Router       /v1/api/products [get]
func (h *ProductHandler) GetAll(c echo.Context) error {
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

	products, total, err := h.service.GetAll(page, pageSize, searchTerm)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := contract.PaginatedResponse[model.Product]{
		Items:      products,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return c.JSON(http.StatusOK, contract.APIResponse[contract.PaginatedResponse[model.Product]]{
		Success: true,
		Data:    resp,
	})
}

// Create godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided information. Category must exist and be valid. Request body should only contain category_id, not the full category object.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.ProductCreateRequest  true  "Product data"
// @Success      201      {object}  model.Product  "Product created successfully"
// @Failure      400      {object}  object{success=bool,error=string}        "Validation error (invalid input or category not found)"
// @Failure      401      {object}  object{success=bool,error=string}        "Unauthorized"
// @Failure      500      {object}  object{success=bool,error=string}        "Internal server error"
// @Security     BearerAuth
// @Router       /v1/api/products [post]
func (h *ProductHandler) Create(c echo.Context) error {
	var req dto.ProductCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	// Validate request
	if err := h.validateCreateRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Convert DTO to model
	product := req.ToProduct()
	if err := h.service.Create(product); err != nil {
		// Check if it's a validation error (category not found, etc.)
		if isValidationError(err) {
			return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
				Success: false,
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, contract.APIResponse[model.Product]{
		Success: true,
		Data:    *product,
	})
}

// GetByID godoc
// @Summary      Get product by ID
// @Description  Retrieve a specific product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID (UUID format)"
// @Success      200  {object}  model.Product
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/products/{id} [get]
func (h *ProductHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}
	if product == nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   "product not found",
		})
	}
	return c.JSON(http.StatusOK, contract.APIResponse[model.Product]{
		Success: true,
		Data:    *product,
	})
}

// Update godoc
// @Summary      Update a product
// @Description  Update an existing product with new information. Category must exist and be valid. Request body should only contain category_id, not the full category object.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Product ID (UUID format)"
// @Param        product  body      dto.ProductUpdateRequest true  "Updated product data"
// @Success      200      {object}  model.Product  "Product updated successfully"
// @Failure      400      {object}  object{success=bool,error=string}        "Validation error (invalid input or category not found)"
// @Failure      401      {object}  object{success=bool,error=string}        "Unauthorized"
// @Failure      404      {object}  object{success=bool,error=string}        "Product not found"
// @Failure      500      {object}  object{success=bool,error=string}        "Internal server error"
// @Security     BearerAuth
// @Router       /v1/api/products/{id} [put]
func (h *ProductHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	var req dto.ProductUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	// Validate request
	if err := h.validateUpdateRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Convert DTO to model
	product := req.ToProduct()

	if err := h.service.Update(id, product); err != nil {
		// Check if it's a validation error (category not found, etc.)
		if isValidationError(err) {
			return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
				Success: false,
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.Product]{
		Success: true,
		Data:    *product,
	})
}

// Delete godoc
// @Summary      Delete a product
// @Description  Delete a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Product ID (UUID format)"
// @Success      204 {string} string "No Content"
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/products/{id} [delete]
func (h *ProductHandler) Delete(c echo.Context) error {
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

// validateCreateRequest validates the product create request
func (h *ProductHandler) validateCreateRequest(req *dto.ProductCreateRequest) error {
	if req.Code == "" {
		return errors.New("product code is required")
	}
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.LargeUnit == "" {
		return errors.New("large unit is required")
	}
	if req.SmallUnit == "" {
		return errors.New("small unit is required")
	}
	if req.ContentPerLargeUnit <= 0 {
		return errors.New("content per large unit must be greater than 0")
	}
	if req.PurchasePrice < 0 {
		return errors.New("purchase price cannot be negative")
	}
	if req.SellingPrice < 0 {
		return errors.New("selling price cannot be negative")
	}
	if req.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	return nil
}

// validateUpdateRequest validates the product update request
func (h *ProductHandler) validateUpdateRequest(req *dto.ProductUpdateRequest) error {
	if req.Code == "" {
		return errors.New("product code is required")
	}
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.LargeUnit == "" {
		return errors.New("large unit is required")
	}
	if req.SmallUnit == "" {
		return errors.New("small unit is required")
	}
	if req.ContentPerLargeUnit <= 0 {
		return errors.New("content per large unit must be greater than 0")
	}
	if req.PurchasePrice < 0 {
		return errors.New("purchase price cannot be negative")
	}
	if req.SellingPrice < 0 {
		return errors.New("selling price cannot be negative")
	}
	if req.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	return nil
}

// isValidationError checks if the error is a validation error
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	validationErrors := []string{
		"is required",
		"cannot be negative",
		"must be greater than 0",
		"not found",
		"invalid",
	}

	for _, validationErr := range validationErrors {
		if strings.Contains(errMsg, validationErr) {
			return true
		}
	}
	return false
}
