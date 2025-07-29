package repository

import (
	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/shared/utils"
	"gorm.io/gorm"
)

type OfficeRepository interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Office, int64, error)
	GetByID(id string) (*model.Office, error)
	Create(office *model.Office) error
	Update(office *model.Office) error
	Delete(id string) error
}

type officeRepository struct {
	db *gorm.DB
}

func NewOfficeRepository(db *gorm.DB) OfficeRepository {
	return &officeRepository{db: db}
}

func (r *officeRepository) GetAll(page, pageSize int, searchTerm string) ([]model.Office, int64, error) {
	var offices []model.Office
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := r.db.Model(&model.Office{})

	// Use reusable search builder for multiple fields
	fields := []string{"code", "name", "address", "city"}
	query = utils.BuildSearchQuery(query, searchTerm, fields)

	// Get total count of matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results with preloaded branches
	if err := query.Preload("Branches").
		Limit(pageSize).
		Offset(offset).
		Find(&offices).Error; err != nil {
		return nil, 0, err
	}

	return offices, total, nil
}

func (r *officeRepository) GetByID(id string) (*model.Office, error) {
	var office model.Office
	if err := r.db.Preload("Branches").First(&office, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &office, nil
}

func (r *officeRepository) Create(office *model.Office) error {
	return r.db.Create(office).Error
}

func (r *officeRepository) Update(office *model.Office) error {
	return r.db.Save(office).Error
}

func (r *officeRepository) Delete(id string) error {
	return r.db.Delete(&model.Office{}, "id = ?", id).Error
}
