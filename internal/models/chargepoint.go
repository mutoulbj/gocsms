package models

import (
	"time"
	
	"github.com/uptrace/bun"
)

type ChargePoint struct {
	bun.BaseModel `bun:"table:charge_points,alias:cp"`
	ID            string    `bun:",pk" json:"id"`
	SerialNumber  string    `bun:"serial_number,notnull" json:"serial_number"`
	Status        string    `bun:"status,notnull" json:"status"`
	LastHeartbeat time.Time `bun:"last_heartbeat,notnull" json:"last_heartbeat"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}
