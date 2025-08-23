package repository

import (
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryProductRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.CategoryProduct, int64, error)
	GetByID(id uuid.UUID) (*model.CategoryProduct, error)
	Create(categoryProduct *model.CategoryProduct) error
	Update(id uuid.UUID, categoryProduct *model.CategoryProduct) error
	Delete(id uuid.UUID) error
	GetByParentID(parentID uuid.UUID) ([]model.CategoryProduct, error)
	GetRootCategories() ([]model.CategoryProduct, error)
}

type categoryProductRepository struct {
	db *gorm.DB
}

func NewCategoryProductRepository(db *gorm.DB) CategoryProductRepository {
	return &categoryProductRepository{db: db}
}

func (r *categoryProductRepository) GetAll(page, pageSize int, searchTerm string) ([]model.CategoryProduct, int64, error) {
	var categories []model.CategoryProduct
	var total int64

	query := r.db.Model(&model.CategoryProduct{})

	if searchTerm != "" {
		query = query.Where("name ILIKE ?", "%"+searchTerm+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryProductRepository) GetByID(id uuid.UUID) (*model.CategoryProduct, error) {
	var category model.CategoryProduct
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryProductRepository) Create(categoryProduct *model.CategoryProduct) error {
	return r.db.Create(categoryProduct).Error
}

func (r *categoryProductRepository) Update(id uuid.UUID, categoryProduct *model.CategoryProduct) error {
	return r.db.Where("id = ?", id).Updates(categoryProduct).Error
}

func (r *categoryProductRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.CategoryProduct{}).Error
}

func (r *categoryProductRepository) GetByParentID(parentID uuid.UUID) ([]model.CategoryProduct, error) {
	var categories []model.CategoryProduct
	if err := r.db.Where("parent_id = ?", parentID).Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryProductRepository) GetRootCategories() ([]model.CategoryProduct, error) {
	var categories []model.CategoryProduct
	if err := r.db.Where("parent_id IS NULL OR parent_id = ?", uuid.Nil).Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
