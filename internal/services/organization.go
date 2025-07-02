package services

import (
	"github.com/mutoulbj/gocsms/internal/repository"
	"github.com/sirupsen/logrus"
)

type OrganizationService struct {
	repo *repository.OrganizationRepository
	log  *logrus.Logger
}

func NewOrganizationService(repo *repository.OrganizationRepository, log *logrus.Logger) *OrganizationService {
	return &OrganizationService{
		repo: repo,
		log:  log,
	}
}

// Create creates a new organization
func (s *OrganizationService) Create(org *repository.Organization) (*repository.Organization, error) {
	s.log.Infof("Creating organization: %s", org.Name)
	return s.repo.Create(org)
}

// GetByID retrieves an organization by its ID
func (s *OrganizationService) GetByID(id uint) (*repository.Organization, error) {
	s.log.Infof("Fetching organization with ID: %d", id)
	return s.repo.GetByID(id)
}

// GetAll retrieves all organizations
func (s *OrganizationService) GetAll(page, pageSize int) ([]*repository.Organization, int64, error) {
	s.log.Info("Fetching all organizations")
	return s.repo.GetAll(page, pageSize)
}

// Update updates an organization
func (s *OrganizationService) Update(org *repository.Organization) (*repository.Organization, error) {
	s.log.Infof("Updating organization with ID: %d", org.ID)
	return s.repo.Update(org)
}

// Delete deletes an organization
func (s *OrganizationService) Delete(id uint) error {
	s.log.Infof("Deleting organization with ID: %d", id)
	return s.repo.Delete(id)
}

// GetByCode retrieves an organization by its code
func (s *OrganizationService) GetByCode(code string) (*repository.Organization, error) {
	s.log.Infof("Fetching organization with code: %s", code)
	return s.repo.GetByCode(code)
}
