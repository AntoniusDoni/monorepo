package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
	"github.com/google/uuid"
)

type UnitProductService interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.UnitProduct, int64, error)
	GetByID(id uuid.UUID) (*model.UnitProduct, error)
	Create(unitProduct *model.UnitProduct) error
	Update(id uuid.UUID, unitProduct *model.UnitProduct) error
	Delete(id uuid.UUID) error
}

type unitProductService struct {
	repo repository.UnitProductRepository
}

func NewUnitProductService(repo repository.UnitProductRepository) UnitProductService {
	return &unitProductService{repo: repo}
}

func (s *unitProductService) GetAll(page, pageSize int, searchTerm string) ([]model.UnitProduct, int64, error) {
	return s.repo.GetAll(page, pageSize, searchTerm)
}

func (s *unitProductService) GetByID(id uuid.UUID) (*model.UnitProduct, error) {
	return s.repo.GetByID(id)
}

func (s *unitProductService) Create(unitProduct *model.UnitProduct) error {
	return s.repo.Create(unitProduct)
}

func (s *unitProductService) Update(id uuid.UUID, unitProduct *model.UnitProduct) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("unit product not found")
	}
	unitProduct.ID = existing.ID // Ensure the ID is set for update
	return s.repo.Update(unitProduct)
}

func (s *unitProductService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}