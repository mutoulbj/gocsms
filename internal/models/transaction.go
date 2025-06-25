package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel  `bun:"table:transactions,alias:tx"`
	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	ChargePointID  uuid.UUID `bun:"charge_point_id,type:uuid,notnull" json:"charge_point_id"`
	ConnectorID    uuid.UUID `bun:"connector_id,type:uuid,notnull" json:"connector_id"`
	TransactionID  int       `bun:"transaction_id,notnull" json:"transaction_id"`
	UserID         uuid.UUID `bun:"user_id,type:uuid,nullzero" json:"user_id"`
	StartTime      time.Time `bun:"start_time,notnull" json:"start_time"`
	StopTime       time.Time `bun:"stop_time,nullzero" json:"stop_time"`
	MeterStart     float64   `bun:"meter_start,notnull" json:"meter_start"`
	MeterStop      float64   `bun:"meter_stop,nullzero" json:"meter_stop"`
	TotalEnergyKwh float64   `bun:"total_energy_kwh,notnull" json:"total_energy_kwh"`
	StopReason     string    `bun:"stop_reason,nullzero" json:"stop_reason"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}
