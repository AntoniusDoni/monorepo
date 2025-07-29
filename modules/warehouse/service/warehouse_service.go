package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
)

type WarehouseService interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Warehouse, int64, error)
	GetByID(id uint) (*model.Warehouse, error)
	Create(warehouse *model.Warehouse) error
	Update(id uint, warehouse *model.Warehouse) error
	Delete(id uint) error
}

type warehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

func (s *warehouseService) GetAll(page, pageSize int, searchTerm string) ([]model.Warehouse, int64, error) {
	return s.repo.GetAll(page, pageSize, searchTerm)
}

func (s *warehouseService) GetByID(id uint) (*model.Warehouse, error) {
	return s.repo.GetByID(id)
}

func (s *warehouseService) Create(warehouse *model.Warehouse) error {
	return s.repo.Create(warehouse)
}

func (s *warehouseService) Update(id uint, warehouse *model.Warehouse) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("warehouse not found")
	}
	warehouse.ID = existing.ID // Ensure the ID is set for update
	return s.repo.Update(warehouse)
}

func (s *warehouseService) Delete(id uint) error {
	return s.repo.Delete(id)
}
