package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MeterValue struct {
	bun.BaseModel `bun:"table:meter_values,alias:mv"`

	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ChargePointID uuid.UUID `bun:"charge_point_id,type:uuid,notnull" json:"charge_point_id"`
	ConnectorID   uuid.UUID `bun:"connector_id,type:uuid,notnull" json:"connector_id"`
	Timestamp     time.Time `bun:"timestamp,notnull" json:"timestamp"`    // Unix timestamp in seconds
	Measurand     string    `bun:"measurand,notnull" json:"measurand"`    // e.g., "Energy.Active.Import.Register"
	Value         float64   `bun:"value,notnull" json:"value"`            // Value in kWh
	Unit          string    `bun:"unit,notnull,default:'Wh'" json:"unit"` // e.g., "kWh", "Wh"
	Context       string    `bun:"context,nullzero" json:"context"`       // Optional context for the meter value
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
}
