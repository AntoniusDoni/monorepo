package handler

import (
	"github.com/antoniusDoni/monorepo/modules/warehouse/dto"
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/service"
	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CategoryProductHandler struct {
	service service.CategoryProductService
}

func NewCategoryProductHandler(service service.CategoryProductService) *CategoryProductHandler {
	return &CategoryProductHandler{service: service}
}

func (h *CategoryProductHandler) RegisterRoutes(g *echo.Group) {
	cpg := g.Group("/category-products")
	cpg.GET("", h.GetAll)
	cpg.POST("", h.Create)
	cpg.GET("/:id", h.GetByID)
	cpg.PUT("/:id", h.Update)
	cpg.DELETE("/:id", h.Delete)
	cpg.GET("/tree", h.GetCategoryTree)
	cpg.GET("/root", h.GetRootCategories)
	cpg.GET("/parent/:parentId", h.GetByParentID)
}

// GetAll godoc
// @Summary      Get list of category products
// @Description  Retrieves paginated category products optionally filtered by search term
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number (default: 1)"
// @Param        pageSize   query     int     false  "Page size (default: 10, max: 100)"
// @Param        searchTerm query     string  false  "Search term to filter categories by name"
// @Success      200        {object}  object
// @Failure      400        {object}  object
// @Failure      401        {object}  object
// @Failure      500        {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products [get]
func (h *CategoryProductHandler) GetAll(c echo.Context) error {
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

	categories, total, err := h.service.GetAll(page, pageSize, searchTerm)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	resp := contract.PaginatedResponse[model.CategoryProduct]{
		Items:      categories,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}

	return c.JSON(http.StatusOK, contract.APIResponse[contract.PaginatedResponse[model.CategoryProduct]]{
		Success: true,
		Data:    resp,
	})
}

// Create godoc
// @Summary      Create a new category product
// @Description  Create a new category product with the provided information. Can be a root category or child category.
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        categoryProduct  body      dto.CategoryProductCreateRequest  true  "Category Product data"
// @Success      201              {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products [post]
func (h *CategoryProductHandler) Create(c echo.Context) error {
	var req dto.CategoryProductCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	categoryProduct := model.CategoryProduct{
		Name:      req.Name,
		UpdatedAt: nil,
	}

	if req.ParentID != uuid.Nil { // only set if not empty
		categoryProduct.ParentID = &req.ParentID
	}

	if err := h.service.Create(&categoryProduct); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, contract.APIResponse[model.CategoryProduct]{
		Success: true,
		Data:    categoryProduct,
	})
}

// GetByID godoc
// @Summary      Get category product by ID
// @Description  Retrieve a specific category product by its ID
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category Product ID (UUID format)"
// @Success      200  {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/{id} [get]
func (h *CategoryProductHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	categoryProduct, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}
	if categoryProduct == nil {
		return c.JSON(http.StatusNotFound, contract.APIResponse[any]{
			Success: false,
			Error:   "category product not found",
		})
	}
	return c.JSON(http.StatusOK, contract.APIResponse[model.CategoryProduct]{
		Success: true,
		Data:    *categoryProduct,
	})
}

// Update godoc
// @Summary      Update a category product
// @Description  Update an existing category product with new information. Validates parent-child relationships and prevents circular references.
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        id              path      string                true  "Category Product ID (UUID format)"
// @Param        categoryProduct body      model.CategoryProduct true  "Updated category product data"
// @Success      200             {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/{id} [put]
func (h *CategoryProductHandler) Update(c echo.Context) error {
	// Parse ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	// Bind request
	var req dto.CategoryProductCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid request body",
		})
	}

	// Map to model
	categoryProduct := model.CategoryProduct{
		Name:      req.Name,
		ParentID:  nil, // default to nil
		CreatedAt: nil,
	}

	// Set ParentID only if provided
	if req.ParentID != uuid.Nil {
		categoryProduct.ParentID = &req.ParentID
	}

	// Call service
	if err := h.service.Update(id, &categoryProduct); err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[model.CategoryProduct]{
		Success: true,
		Data:    categoryProduct,
	})
}

// Delete godoc
// @Summary      Delete a category product
// @Description  Delete a category product by its ID. Cannot delete categories that have child categories.
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Category Product ID (UUID format)"
// @Success      204 {string} string "No Content"
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/{id} [delete]
func (h *CategoryProductHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid id format",
		})
	}

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "cannot delete category with child categories" {
			return c.JSON(http.StatusConflict, contract.APIResponse[any]{
				Success: false,
				Error:   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetCategoryTree godoc
// @Summary      Get category tree structure
// @Description  Retrieve all categories organized in a hierarchical tree structure
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/tree [get]
func (h *CategoryProductHandler) GetCategoryTree(c echo.Context) error {
	tree, err := h.service.GetCategoryTree()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[[]service.CategoryTreeNode]{
		Success: true,
		Data:    tree,
	})
}

// GetRootCategories godoc
// @Summary      Get root categories
// @Description  Retrieve all root categories (categories without parent)
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/root [get]
func (h *CategoryProductHandler) GetRootCategories(c echo.Context) error {
	categories, err := h.service.GetRootCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[[]model.CategoryProduct]{
		Success: true,
		Data:    categories,
	})
}

// GetByParentID godoc
// @Summary      Get categories by parent ID
// @Description  Retrieve all child categories of a specific parent category
// @Tags         category-products
// @Accept       json
// @Produce      json
// @Param        parentId  path      string  true  "Parent Category ID (UUID format)"
// @Success      200  {object}  model.CategoryProduct
// @Failure      400  {object}  object
// @Failure      401  {object}  object
// @Failure      404  {object}  object
// @Failure      500  {object}  object
// @Security     BearerAuth
// @Router       /v1/api/category-products/parent/{parentId} [get]
func (h *CategoryProductHandler) GetByParentID(c echo.Context) error {
	parentIdStr := c.Param("parentId")
	parentId, err := uuid.Parse(parentIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, contract.APIResponse[any]{
			Success: false,
			Error:   "invalid parent id format",
		})
	}

	categories, err := h.service.GetByParentID(parentId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, contract.APIResponse[any]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, contract.APIResponse[[]model.CategoryProduct]{
		Success: true,
		Data:    categories,
	})
}
