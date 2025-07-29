package repository

import (
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/shared/utils"
	"gorm.io/gorm"
)

type BranchRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Branch, int64, error)
	GetByID(id string) (*model.Branch, error)
	Create(branch *model.Branch) error
	Update(branch *model.Branch) error
	Delete(id string) error
	GetByOfficeID(officeID string) ([]model.Branch, error)
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) GetAll(page, pageSize int, searchTerm string) ([]model.Branch, int64, error) {
	var branches []model.Branch
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := r.db.Model(&model.Branch{})

	// Optional: add reusable search helper for fields like code, name, city
	fields := []string{"code", "name", "address", "city"}
	query = utils.BuildSearchQuery(query, searchTerm, fields)

	// Count total records matching filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query with preload and pagination
	if err := query.Preload("Office").
		Preload("Warehouses").
		Limit(pageSize).
		Offset(offset).
		Find(&branches).Error; err != nil {
		return nil, 0, err
	}

	return branches, total, nil
}

func (r *branchRepository) GetByID(id string) (*model.Branch, error) {
	var branch model.Branch
	if err := r.db.Preload("Office").Preload("Warehouses").First(&branch, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) Create(branch *model.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) Update(branch *model.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id string) error {
	return r.db.Delete(&model.Branch{}, "id = ?", id).Error
}

func (r *branchRepository) GetByOfficeID(officeID string) ([]model.Branch, error) {
	var branches []model.Branch
	if err := r.db.Where("office_id = ?", officeID).Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}
