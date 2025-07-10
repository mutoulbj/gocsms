package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/dto"
	"github.com/mutoulbj/gocsms/internal/models"
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

// Create from request
func (s *OrganizationService) CreateFromRequest(
	ctx context.Context,
	req *dto.OrganizationCreateRequest,
) (*models.Organization, error) {
	s.log.Infof("Creating organization from request: %s", req.Name)
	org := &models.Organization{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	created, err := s.repo.Create(ctx, org)
	if err != nil {
		s.log.WithError(err).Error("Failed to create organization from request")
		return nil, err
	}
	return created, nil
}

// Create creates a new organization
func (s *OrganizationService) Create(
	ctx context.Context,
	org *models.Organization,
) (*dto.OrganizationResponse, error) {
	s.log.Infof("Creating organization: %s", org.Name)
	created, err := s.repo.Create(ctx, org)
	if err != nil {
		s.log.WithError(err).Error("Failed to create organization")
		return nil, err
	}
	return dto.ToOrganizationResponse(created), nil
}

// GetByID retrieves an organization by its ID
func (s *OrganizationService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.OrganizationResponse, error) {
	s.log.Infof("Fetching organization with ID: %s", id)
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("Failed to fetch organization")
		return nil, err
	}
	return dto.ToOrganizationResponse(result), nil
}

// GetAll retrieves all organizations
func (s *OrganizationService) GetAll(ctx context.Context, page, pageSize int) ([]*dto.OrganizationResponse, int64, error) {
	s.log.Info("Fetching all organizations")
	orgs, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.OrganizationResponse
	for _, org := range orgs {
		responses = append(responses, dto.ToOrganizationResponse(org))
	}
	return responses, total, nil
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
