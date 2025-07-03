package services

import (
	"context"

	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/google/uuid"
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

// Create from request
func (s *OrganizationService) CreateFromRequest(ctx context.Context, req *models.OrganizationCreateRequest) (*models.Organization, error) {
	s.log.Infof("Creating organization from request: %s", req.Name)
	org := &models.Organization{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	err := s.repo.Create(ctx, org)
	if err != nil {
		s.log.WithError(err).Error("Failed to create organization from request")
		return nil, err
	}
	return org, nil
}

// Create creates a new organization
func (s *OrganizationService) Create(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	s.log.Infof("Creating organization: %s", org.Name)
	err := s.repo.Create(ctx, org)
	if err != nil {
		s.log.WithError(err).Error("Failed to create organization")
		return nil, err
	}
	return org, nil
}

// GetByID retrieves an organization by its ID
func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	s.log.Infof("Fetching organization with ID: %s", id)
	return s.repo.GetByID(ctx, id)
}

// GetAll retrieves all organizations
func (s *OrganizationService) GetAll(ctx context.Context, page, pageSize int) ([]*models.Organization, int64, error) {
	s.log.Info("Fetching all organizations")
	return s.repo.List(ctx, page, pageSize)
}

// Update updates an organization
func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, org *models.Organization) (*models.Organization, error) {
	s.log.Infof("Updating organization with ID: %d", org.ID)
	err := s.repo.Update(ctx, id, org)
	if err != nil {
		s.log.WithError(err).Error("Failed to update organization")
		return nil, err
	}
	return org, nil
}

// Delete deletes an organization
func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	s.log.Infof("Deleting organization with ID: %s", id)
	return s.repo.Delete(ctx, id)
}

// GetByCode retrieves an organization by its code
func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	s.log.Infof("Fetching organization with slug: %s", slug)
	return s.repo.GetBySlug(ctx, slug)
}
