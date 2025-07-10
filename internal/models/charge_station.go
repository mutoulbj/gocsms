package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ChargeStation struct {
	bun.BaseModel `bun:"table:charge_stations,alias:cs"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Address   string    `bun:"address,notnull" json:"address"`
	City      string    `bun:"city,notnull" json:"city"`
	State     string    `bun:"state,notnull" json:"state"`
	Country   string    `bun:"country,notnull" json:"country"`
	Latitude  float64   `bun:"latitude,notnull" json:"latitude"`
	Longitude float64   `bun:"longitude,notnull" json:"longitude"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	OrganizationID uuid.UUID     `bun:"organization_id,type:uuid,notnull" json:"organization_id"` // Foreign key to Organization
	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id" json:"organization,omitempty"`
}

func (cs *ChargeStation) BeforeInsert() error {
	cs.ID = uuid.New()
	cs.CreatedAt = time.Now()
	cs.UpdatedAt = time.Now()
	return nil
}

func (cs *ChargeStation) BeforeUpdate() error {
	cs.UpdatedAt = time.Now()
	return nil
}
