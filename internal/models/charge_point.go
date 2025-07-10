package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/enums"

	"github.com/uptrace/bun"
)

type ChargePoint struct {
	bun.BaseModel      `bun:"table:charge_points,alias:cp"`
	ID                 uuid.UUID                           `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name               string                              `bun:"name,notnull" json:"name"`
	Code               string                              `bun:"code,notnull,unique" json:"code"`
	SerialNumber       string                              `bun:"serial_number,notnull" json:"serial_number"`
	Status             enums.ChargePointStatus             `bun:"status,type:VARCHAR(20),default:UNKNOWN,notnull" json:"status"`
	LastHeartbeat      time.Time                           `bun:"last_heartbeat,notnull" json:"last_heartbeat"`
	OcppVersion        string                              `bun:"ocpp_version,notnull" json:"ocpp_version"`
	RegistrationStatus enums.ChargePointRegistrationStatus `bun:"registration_status,notnull" json:"registration_status"`
	RegisteredAt       time.Time                           `bun:"registered_at" json:"registered_at"`
	CreatedAt          time.Time                           `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time                           `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	Model              string                              `bun:"model" json:"model"`
	Vendor             string                              `bun:"vendor" json:"vendor"`
	Connected          bool                                `bun:"connected,notnull,default:false" json:"connected"`
	ChargeStationId    uuid.UUID                           `bun:"charge_station_id,notnull" json:"charge_station_id"`
	Connectors         []*Connector                        `bun:"rel:has-many,join:id=charge_point_id" json:"connectors,omitempty"`
	ChargeStation      *ChargeStation                      `bun:"rel:belongs-to,join:charge_station_id=id" json:"charge_station,omitempty"`
}

func (c *ChargePoint) BeforeInsert() error {
	c.ID = uuid.New()
	c.RegistrationStatus = enums.ChargePointRegistrationStatusPending
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	if c.Status == "" {
		c.Status = enums.ChargePointStatusUnknown
	}
	return nil
}

func (c *ChargePoint) BeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}
