package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ChargeStation struct {
	bun.BaseModel `bun:"table:charge_stations,alias:cs"`

	ID              uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Name            string    `bun:"name,notnull" json:"name"`
	Identity        string    `bun:"identity,notnull" json:"identity"`
	Vendor          string    `bun:"vendor,notnull" json:"vendor"`
	Model           string    `bun:"model,notnull" json:"model"`
	SerialNumber    string    `bun:"serial_number,notnull" json:"serial_number"`
	FirewareVersion string    `bun:"fireware_version,notnull" json:"fireware_version"`
	// Location  Infomation
	AddressLine1 string  `bun:"address_line1,notnull" json:"address_line1"`
	AddressLine2 string  `bun:"address_line2,nullzero" json:"address_line2"`
	City         string  `bun:"city,notnull" json:"city"`
	State        string  `bun:"state,notnull" json:"state"`
	Country      string  `bun:"country,notnull" json:"country"`
	Latitude     float64 `bun:"latitude,notnull" json:"latitude"`
	Longitude    float64 `bun:"longitude,notnull" json:"longitude"`
	CreatedAt    string  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    string  `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	OrganizationID uuid.UUID `bun:"organization_id,type:uuid,notnull" json:"organization_id"` // Foreign key to Organization

	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id" json:"organization,omitempty"`
}
