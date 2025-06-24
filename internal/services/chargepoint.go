package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/repository"
)

type ChargePointService struct {
	repo *repository.ChargePointRepository
	log  *logrus.Logger
}

func GocsmsChargePointService(repo *repository.ChargePointRepository, log *logrus.Logger) *ChargePointService {
	return &ChargePointService{repo: repo, log: log}
}

func (s *ChargePointService) Register(ctx context.Context, cp *models.ChargePoint) error {
	return s.repo.Create(ctx, cp)
}

func (s *ChargePointService) GetByID(ctx context.Context, id string) (*models.ChargePoint, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ChargePointService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
