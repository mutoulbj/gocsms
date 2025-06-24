package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"

	"github.com/mutoulbj/gocsms/internal/models"
)

type ChargePointRepository struct {
	db    *bun.DB
	redis *redis.Client
	log   *logrus.Logger
}

func GocsmsChargePointRepository(db *bun.DB, redis *redis.Client, log *logrus.Logger) *ChargePointRepository {
	return &ChargePointRepository{
		db:    db,
		redis: redis,
		log:   log,
	}
}

func (r *ChargePointRepository) Create(ctx context.Context, cp *models.ChargePoint) error {
	_, err := r.db.NewInsert().Model(cp).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).Error("failed to create charge point")
		return err
	}
	return r.cacheChargePoint(ctx, cp)
}

func (r *ChargePointRepository) GetByID(ctx context.Context, id string) (*models.ChargePoint, error) {
	// try cache first
	cp, err := r.getFromCache(ctx, id)
	if err == nil && cp != nil {
		return cp, nil
	}

	// if not found in cache, try db
	cp = &models.ChargePoint{}
	err = r.db.NewSelect().Model(cp).Where("id = ?", id).Scan(ctx)
	if err != nil {
		r.log.Error("failed to get charge point from db")
		return nil, err
	}
	return cp, r.cacheChargePoint(ctx, cp)
}

func (r *ChargePointRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.NewUpdate().
		Model((*models.ChargePoint)(nil)).
		Set("status = ?, updated_at = ?", status, time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		r.log.Error("failed to update charge point status: ", err)
		return err
	}
	return r.invalidateCache(ctx, id.String())
}

func (r *ChargePointRepository) cacheChargePoint(ctx context.Context, cp *models.ChargePoint) error {
	data, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, "chargepoint:"+cp.ID.String(), data, 5*time.Minute).Err()
}

func (r *ChargePointRepository) getFromCache(ctx context.Context, id string) (*models.ChargePoint, error) {
	data, err := r.redis.Get(ctx, "chargepoint:"+id).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cp models.ChargePoint
	return &cp, json.Unmarshal(data, &cp)
}

func (r *ChargePointRepository) invalidateCache(ctx context.Context, id string) error {
	return r.redis.Del(ctx, "chargepoint:"+id).Err()
}
