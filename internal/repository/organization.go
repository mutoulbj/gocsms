package repository

import (
	"context"

	"gocsms/internal/domain"

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

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	_, err := r.db.NewInsert().Model(org).Exec(ctx)
	if err != nil {
		if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() {
			return domain.ErrDuplicateOrganization
		}
		r.log.WithError(err).Error("Failed to create organization")
		return err
	}
	return nil
}

// GetByID retrieves an organization by its ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id int64) (*domain.Organization, error) {
	org := new(domain.Organization)
	err := r.db.NewSelect().Model(org).Where("id = ?", id).Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to get organization by ID")
		return nil, err
	}
	return org, nil
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	_, err := r.db.NewUpdate().Model(org).WherePK().Exec(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to update organization")
		return err
	}
	return nil
}

// Delete deletes an organization by its ID
func (r *OrganizationRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Organization)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to delete organization")
		return err
	}
	return nil
}

// List returns all organizations
func (r *OrganizationRepository) List(ctx context.Context) ([]*domain.Organization, error) {
	var orgs []*domain.Organization
	err := r.db.NewSelect().Model(&orgs).Scan(ctx)
	if err != nil {
		r.log.WithError(err).Error("Failed to list organizations")
		return nil, err
	}
	return orgs, nil
}
