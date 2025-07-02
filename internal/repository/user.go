package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db  *bun.DB
	log *logrus.Logger
}

func NewUserRepository(db *bun.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		r.log.Error("Failed to get user by ID: ", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		r.log.Error("Failed to get user by username: ", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		r.log.Error("Failed to create user: ", err)
		return err
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		Column("first_name", "last_name").
		Where("id = ?", user.ID).
		Exec(ctx)
	if err != nil {
		r.log.Error("Failed to update user: ", err)
		return err
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	total, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Count(ctx)
	if err != nil {
		r.log.Error("Failed to count users: ", err)
		return nil, 0, err
	}
	err = r.db.NewSelect().
		Model(&users).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		r.log.Error("Failed to list users: ", err)
		return nil, 0, err
	}
	return users, int64(total), nil
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("username = ?", username).
		Count(ctx)
	if err != nil {
		r.log.Error("Failed to check if username exists: ", err)
		return false, err
	}
	return count > 0, nil
}
