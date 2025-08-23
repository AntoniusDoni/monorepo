package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	"github.com/google/uuid"
)

type CategoryProductService interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.CategoryProduct, int64, error)
	GetByID(id uuid.UUID) (*model.CategoryProduct, error)
	Create(categoryProduct *model.CategoryProduct) error
	Update(id uuid.UUID, categoryProduct *model.CategoryProduct) error
	Delete(id uuid.UUID) error
	GetByParentID(parentID uuid.UUID) ([]model.CategoryProduct, error)
	GetRootCategories() ([]model.CategoryProduct, error)
	GetCategoryTree() ([]CategoryTreeNode, error)
}

type CategoryTreeNode struct {
	ID       uuid.UUID          `json:"id"`
	Name     string             `json:"name"`
	ParentID *uuid.UUID         `json:"parent_id"`
	Children []CategoryTreeNode `json:"children,omitempty"`
}

type categoryProductService struct {
	repo repository.CategoryProductRepository
}

func NewCategoryProductService(repo repository.CategoryProductRepository) CategoryProductService {
	return &categoryProductService{repo: repo}
}

func (s *categoryProductService) GetAll(page, pageSize int, searchTerm string) ([]model.CategoryProduct, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limit maximum page size
	}

	return s.repo.GetAll(page, pageSize, searchTerm)
}

func (s *categoryProductService) GetByID(id uuid.UUID) (*model.CategoryProduct, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}
	return s.repo.GetByID(id)
}

func (s *categoryProductService) Create(categoryProduct *model.CategoryProduct) error {
	if categoryProduct == nil {
		return errors.New("category product cannot be nil")
	}
	if categoryProduct.Name == "" {
		return errors.New("category name is required")
	}
	if categoryProduct.ParentID != nil {
		parent, err := s.repo.GetByID(*categoryProduct.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent category not found")
		}

		// Prevent circular reference (category cannot be its own parent)
		if err := s.validateNoCircularReference(categoryProduct.ID, categoryProduct.ParentID); err != nil {
			return err
		}
	}

	return s.repo.Create(categoryProduct)
}

func (s *categoryProductService) Update(id uuid.UUID, categoryProduct *model.CategoryProduct) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}
	if categoryProduct == nil {
		return errors.New("category product cannot be nil")
	}
	if categoryProduct.Name == "" {
		return errors.New("category name is required")
	}

	// Check if category exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}

	// Validate parent category exists if ParentID is provided and different from current ID
	if categoryProduct.ParentID != nil && *categoryProduct.ParentID != id {
		parent, err := s.repo.GetByID(*categoryProduct.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent category not found")
		}

		// Prevent circular reference
		if err := s.validateNoCircularReference(id, categoryProduct.ParentID); err != nil {
			return err
		}
	}

	return s.repo.Update(id, categoryProduct)
}

func (s *categoryProductService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}

	// Check if category exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}

	// Check if category has children
	children, err := s.repo.GetByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete category with child categories")
	}

	return s.repo.Delete(id)
}

func (s *categoryProductService) GetByParentID(parentID uuid.UUID) ([]model.CategoryProduct, error) {
	return s.repo.GetByParentID(parentID)
}

func (s *categoryProductService) GetRootCategories() ([]model.CategoryProduct, error) {
	return s.repo.GetRootCategories()
}

func (s *categoryProductService) GetCategoryTree() ([]CategoryTreeNode, error) {
	// Get all categories
	allCategories, _, err := s.repo.GetAll(1, 1000, "") // Get all categories
	if err != nil {
		return nil, err
	}

	// Build category map for quick lookup
	categoryMap := make(map[uuid.UUID]*CategoryTreeNode)
	for _, cat := range allCategories {
		categoryMap[cat.ID] = &CategoryTreeNode{
			ID:       cat.ID,
			Name:     cat.Name,
			ParentID: cat.ParentID,
			Children: []CategoryTreeNode{},
		}
	}

	// Build tree structure
	var rootNodes []CategoryTreeNode
	for _, cat := range allCategories {
		node := categoryMap[cat.ID]

		if cat.ParentID == nil {
			// Root category (no parent, NULL in DB)
			rootNodes = append(rootNodes, *node)
		} else {
			// Child category
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	return rootNodes, nil
}

func (s *categoryProductService) validateNoCircularReference(categoryID uuid.UUID, parentID *uuid.UUID) error {
	// If no parent, nothing to validate
	if parentID == nil {
		return nil
	}

	// Traverse up the parent chain to check for circular reference
	currentID := parentID
	visited := make(map[uuid.UUID]bool)

	for currentID != nil {
		if visited[*currentID] {
			return errors.New("circular reference detected")
		}
		if *currentID == categoryID {
			return errors.New("cannot set category as its own parent")
		}

		visited[*currentID] = true

		parent, err := s.repo.GetByID(*currentID)
		if err != nil {
			return err
		}
		if parent == nil || parent.ParentID == nil {
			break
		}

		currentID = parent.ParentID
	}

	return nil
}
