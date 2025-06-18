package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel   `bun:"table:transactions,alias:tx"`
	ID              int64     `bun:",pk,autoincrement" json:"id"`
	ChargePointID   string    `bun:"charge_point_id,notnull" json:"charge_point_id"`
	ConnectorID     int       `bun:"connector_id,notnull" json:"connector_id"`
	TransactionID   int       `bun:"transaction_id,notnull" json:"transaction_id"`
	StartTime       time.Time `bun:"start_time,notnull" json:"start_time"`
	StopTime        time.Time `bun:"stop_time,nullzero" json:"stop_time"`
	MeterStart      float64   `bun:"meter_start,notnull" json:"meter_start"`
	MeterStop       float64   `bun:"meter_stop,nullzero" json:"meter_stop"`
	CreatedAt       time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}