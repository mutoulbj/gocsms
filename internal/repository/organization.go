package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type OrganizationRepository struct {
	db  *bun.DB
	log *logrus.Logger
}

func NewOrganizationRepository(db *bun.DB, log *logrus.Logger) *OrganizationRepository {
	return &OrganizationRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new organization and returns the created object
func (r *OrganizationRepository) Create(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	err := r.db.NewInsert().
		Model(org).
		Returning("*").
		Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to create organization")
		return nil, err
	}
	return org, nil
}

// GetByID retrieves an organization by its ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org := &models.Organization{}
	err := r.db.NewSelect().
		Model(org).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to get organization by ID")
		return nil, err
	}
	return org, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, id uuid.UUID, org *models.Organization) error {
	_, err := r.db.NewUpdate().
		Model(org).
		Column("name", "slug", "description").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to update organization")
		return err
	}
	return nil
}

// Delete deletes an organization by its ID
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*models.Organization)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to delete organization")
		return err
	}
	return nil
}

// List returns all organizations
func (r *OrganizationRepository) List(ctx context.Context, offset, limit int) ([]*models.Organization, int64, error) {
	var orgs []*models.Organization

	// total count of organizations
	total, err := r.db.NewSelect().
		Model((*models.Organization)(nil)).
		Count(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to count organizations")
		return nil, 0, err
	}

	// get organization list
	err = r.db.NewSelect().
		Model(&orgs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to list organizations")
		return nil, 0, err
	}
	return orgs, int64(total), nil
}

// GetBySlug retrieves an organization by its slug
func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	org := &models.Organization{}
	err := r.db.NewSelect().
		Model(org).
		Where("slug = ?", slug).
		Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to get organization by slug")
		return nil, err
	}
	return org, nil
}
