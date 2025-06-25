// Package models provides data structures and database models for the application.
//
// organization.go defines the Organization model and related logic.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:organizations,alias:o"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name          string    `bun:"name,notnull" json:"name"`
	Slug          string    `bun:"slug,notnull,unique" json:"slug"`
	Description   string    `bun:"description,nullzero" json:"description"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	ChargeStations []*ChargeStation `bun:"rel:has-many,join:id=organization_id" json:"charge_stations,omitempty"`
}
