package repository

import (
	"context"
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Product, int64, error)
	GetByID(id uuid.UUID) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	*repository.Repository
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{Repository: repository.NewRepository(context.Background(), db)}
}

func (r *productRepository) GetAll(page, pageSize int, searchTerm string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	baseQuery := r.DB().Model(&model.Product{}).Preload("Category")
	if searchTerm != "" {
		searchTerm = utils.SanitizeSearchTerm(searchTerm)
		like := "%" + searchTerm + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR code ILIKE ? OR indication ILIKE ?", like, like, like)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.DB().Preload("Category").First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // not found, return nil object and nil error
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.DB().Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.DB().Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.DB().Delete(&model.Product{}, "id = ?", id).Error
}
