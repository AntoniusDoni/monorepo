package service

import (
	"errors"

	"github.com/antoniusDoni/monorepo/modules/warehouse/model"
	"github.com/antoniusDoni/monorepo/modules/warehouse/repository"
)

type OfficeService interface {
	GetAll(page, pageSize int, searchTerm string) ([]model.Office, int64, error)
	GetActiveOffices() ([]model.Office, error)
	GetByID(id string) (*model.Office, error)
	Create(office *model.Office) error
	Update(id string, office *model.Office) error
	Delete(id string) error
}
type officeService struct {
	repo repository.OfficeRepository
}

func NewOfficeService(repo repository.OfficeRepository) OfficeService {
	return &officeService{repo: repo}
}
func (s *officeService) GetAll(page, pageSize int, searchTerm string) ([]model.Office, int64, error) {
	return s.repo.GetAll(page, pageSize, searchTerm)
}

func (s *officeService) GetActiveOffices() ([]model.Office, error) {
	return s.repo.GetActiveOffices()
}

func (s *officeService) GetByID(id string) (*model.Office, error) {
	return s.repo.GetByID(id)
}
func (s *officeService) Create(office *model.Office) error {
	return s.repo.Create(office)
}
func (s *officeService) Update(id string, office *model.Office) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("Office not found")
	}
	office.ID = existing.ID // Ensure the ID is set for update
	return s.repo.Update(office)
}
func (s *officeService) Delete(id string) error {
	return s.repo.Delete(id)
}
