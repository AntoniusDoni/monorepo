package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Product, int64, error)
	GetByID(id uuid.UUID) (*model.Product, error)
	Create(product *model.Product) error
	Update(id uuid.UUID, product *model.Product) error
	Delete(id uuid.UUID) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryProductRepository
}

func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryProductRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) GetAll(page, pageSize int, searchTerm string) ([]model.Product, int64, error) {
	return s.repo.GetAll(page, pageSize, searchTerm)
}

func (s *productService) GetByID(id uuid.UUID) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Create(product *model.Product) error {
	// Validate required fields
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Validate category exists
	if err := s.validateCategoryExists(product.CategoryID); err != nil {
		return err
	}

	return s.repo.Create(product)
}

func (s *productService) Update(id uuid.UUID, product *model.Product) error {
	// Check if product exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	// Validate required fields
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Validate category exists
	if err := s.validateCategoryExists(product.CategoryID); err != nil {
		return err
	}

	product.ID = existing.ID // Ensure the ID is set for update
	return s.repo.Update(product)
}

// validateProduct validates the product fields
func (s *productService) validateProduct(product *model.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}
	if product.Code == "" {
		return errors.New("product code is required")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.LargeUnit == "" {
		return errors.New("large unit is required")
	}
	if product.SmallUnit == "" {
		return errors.New("small unit is required")
	}
	if product.ContentPerLargeUnit <= 0 {
		return errors.New("content per large unit must be greater than 0")
	}
	if product.PurchasePrice < 0 {
		return errors.New("purchase price cannot be negative")
	}
	if product.SellingPrice < 0 {
		return errors.New("selling price cannot be negative")
	}
	if product.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	return nil
}

// validateCategoryExists checks if the category exists
func (s *productService) validateCategoryExists(categoryID uuid.UUID) error {
	if categoryID == uuid.Nil {
		return errors.New("category ID is required")
	}

	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return errors.New("failed to validate category: " + err.Error())
	}
	if category == nil {
		return errors.New("category not found")
	}
	return nil
}

func (s *productService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}