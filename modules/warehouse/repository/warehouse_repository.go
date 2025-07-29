package repository

import (
	"context"
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/utils"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Warehouse, int64, error)
	GetByID(id uint) (*model.Warehouse, error)
	Create(warehouse *model.Warehouse) error
	Update(warehouse *model.Warehouse) error
	Delete(id uint) error
}

type warehouseRepository struct {
	*repository.Repository
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{Repository: repository.NewRepository(context.Background(), db)}
}

func (r *warehouseRepository) GetAll(page, pageSize int, searchTerm string) ([]model.Warehouse, int64, error) {
	var warehouses []model.Warehouse
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	baseQuery := r.DB().Model(&model.Warehouse{})
	if searchTerm != "" {
		searchTerm = utils.SanitizeSearchTerm(searchTerm)
		like := "%" + searchTerm + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR location ILIKE ?", like, like)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Limit(pageSize).Offset(offset).Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}
	return warehouses, total, nil
}

func (r *warehouseRepository) GetByID(id uint) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	err := r.DB().First(&warehouse, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // not found, return nil object and nil error
	}
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) Create(warehouse *model.Warehouse) error {
	return r.DB().Create(warehouse).Error
}

func (r *warehouseRepository) Update(warehouse *model.Warehouse) error {
	return r.DB().Save(warehouse).Error
}

func (r *warehouseRepository) Delete(id uint) error {
	return r.DB().Delete(&model.Warehouse{}, id).Error
}
