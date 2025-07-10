package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/enums"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel  `bun:"table:users,alias:u"`
	ID             uuid.UUID        `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Username       string           `bun:"username,notnull,unique" json:"username"`
	FirstName      string           `bun:"first_name,nullzero" json:"first_name"`
	LastName       string           `bun:"last_name,nullzero" json:"last_name"`
	Email          string           `bun:"email,notnull,unique" json:"email"`
	PhoneNumber    string           `bun:"phone_number,nullzero" json:"phone_number"`
	PasswordHash   string           `bun:"password_hash,notnull" json:"-"`
	CreatedAt      time.Time        `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time        `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	LastLoginAt    time.Time        `bun:"last_login_at,nullzero" json:"last_login_at"`
	Status         enums.UserStatus `bun:"status,notnull,default:'ACTIVE'" json:"status"`
	OrganizationID uuid.UUID        `bun:"organization_id,nullzero" json:"organization_id"`
	Organization   *Organization    `bun:"rel:belongs-to,join:organization_id=id" json:"organization,omitempty"`
}

func (u *User) BeforeInsert(ctx context.Context) error {
	// Set CreatedAt and UpdatedAt before inserting
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if u.Status == "" {
		u.Status = enums.UserStatusActive
	}
	return nil
}

func (u *User) BeforeUpdate(ctx context.Context) error {
	// Update UpdatedAt before updating
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) UserResponse() *User {
	return &User{
		ID:             u.ID,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		Status:         u.Status,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		LastLoginAt:    u.LastLoginAt,
		OrganizationID: u.OrganizationID,
	}
}
