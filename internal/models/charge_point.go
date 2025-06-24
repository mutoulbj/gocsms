package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/enums"

	"github.com/uptrace/bun"
)

type ChargePoint struct {
	bun.BaseModel      `bun:"table:charge_points,alias:cp"`
	ID                 uuid.UUID                           `bun:",pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Name               string                              `bun:"name,notnull" json:"name"`
	Code               string                              `bun:"code,notnull" json:"code"`
	SerialNumber       string                              `bun:"serial_number,notnull" json:"serial_number"`
	Status             enums.ChargePointStatus             `bun:"status,notnull" json:"status"`
	LastHeartbeat      time.Time                           `bun:"last_heartbeat,notnull" json:"last_heartbeat"`
	OcppProtocol       string                              `bun:"ocpp_protocol,notnull" json:"ocpp_protocol"`
	RegistrationStatus enums.ChargePointRegistrationStatus `bun:"registration_status,notnull" json:"registration_status"`
	CreatedAt          time.Time                           `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time                           `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	ChargeStationId uuid.UUID `bun:"charge_station_id,notnull" json:"charge_station_id"`

	Connectors    []*Connector   `bun:"rel:has-many,join:id=charge_point_id" json:"connectors,omitempty"`
	ChargeStation *ChargeStation  `bun:"rel:belongs-to,join:charge_station_id=id" json:"charge_station,omitempty"`
}

type Connector struct {
	bun.BaseModel `bun:"table:connectors,alias:c"`
	ID            uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	ChargePointID uuid.UUID `bun:"charge_point_id,type:uuid,notnull" json:"charge_point_id"`
	Standard      string    `bun:"standard,notnull" json:"standard"`
	Format        string    `bun:"format,notnull" json:"format"`
	PowerType     string    `bun:"power_type,notnull" json:"power_type"`
	MaxVoltage    int       `bun:"max_voltage,notnull" json:"max_voltage"`
	MaxAmperage   int       `bun:"max_amperage,notnull" json:"max_amperage"`
	MaxPower      int       `bun:"max_power,notnull" json:"max_power"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	ChargePoint *ChargePoint `bun:"rel:belongs-to,join:charge_point_id=id" json:"charge_point,omitempty"`
}

type CreateChargePointRequest struct {
	Name               string `json:"name" validate:"required"`
	Code               string `json:"code" validate:"required"`
	Status             string `json:"status" validate:"required,oneof=available unavailable charging paused reserved"`
	RegistrationStatus string `json:"registration_status" validate:"required,oneof=accepted rejected pending"`
}
