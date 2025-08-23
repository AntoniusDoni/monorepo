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

type UnitProductRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.UnitProduct, int64, error)
	GetByID(id uuid.UUID) (*model.UnitProduct, error)
	Create(unitProduct *model.UnitProduct) error
	Update(unitProduct *model.UnitProduct) error
	Delete(id uuid.UUID) error
}

type unitProductRepository struct {
	*repository.Repository
}

func NewUnitProductRepository(db *gorm.DB) UnitProductRepository {
	return &unitProductRepository{Repository: repository.NewRepository(context.Background(), db)}
}

func (r *unitProductRepository) GetAll(page, pageSize int, searchTerm string) ([]model.UnitProduct, int64, error) {
	var unitProducts []model.UnitProduct
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	baseQuery := r.DB().Model(&model.UnitProduct{})
	if searchTerm != "" {
		searchTerm = utils.SanitizeSearchTerm(searchTerm)
		like := "%" + searchTerm + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Limit(pageSize).Offset(offset).Find(&unitProducts).Error; err != nil {
		return nil, 0, err
	}
	return unitProducts, total, nil
}

func (r *unitProductRepository) GetByID(id uuid.UUID) (*model.UnitProduct, error) {
	var unitProduct model.UnitProduct
	err := r.DB().First(&unitProduct, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // not found, return nil object and nil error
	}
	if err != nil {
		return nil, err
	}
	return &unitProduct, nil
}

func (r *unitProductRepository) Create(unitProduct *model.UnitProduct) error {
	return r.DB().Create(unitProduct).Error
}

func (r *unitProductRepository) Update(unitProduct *model.UnitProduct) error {
	return r.DB().Save(unitProduct).Error
}

func (r *unitProductRepository) Delete(id uuid.UUID) error {
	return r.DB().Delete(&model.UnitProduct{}, "id = ?", id).Error
}
