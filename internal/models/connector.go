package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Connector struct {
	bun.BaseModel `bun:"table:connectors,alias:c"`
	ID            uuid.UUID    `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	ChargePointID uuid.UUID    `bun:"charge_point_id,type:uuid,notnull" json:"charge_point_id"`
	ConnectorID   string       `bun:"connector_id,notnull" json:"connector_id"`
	Standard      string       `bun:"standard,notnull" json:"standard"`
	Format        string       `bun:"format,notnull" json:"format"`
	PowerType     string       `bun:"power_type,notnull" json:"power_type"`
	MaxVoltage    int          `bun:"max_voltage,notnull" json:"max_voltage"`
	MaxAmperage   int          `bun:"max_amperage,notnull" json:"max_amperage"`
	MaxPower      int          `bun:"max_power,notnull" json:"max_power"`
	CreatedAt     time.Time    `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time    `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	ChargePoint   *ChargePoint `bun:"rel:belongs-to,join:charge_point_id=id" json:"charge_point,omitempty"`
}

func (c *Connector) BeforeInsert() error {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Connector) BeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}
