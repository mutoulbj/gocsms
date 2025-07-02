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

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	// Perform any necessary actions before appending the model
	switch query.(type) {
	case *bun.InsertQuery:
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		if u.Status == "" {
			u.Status = enums.UserStatusActive
		}
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now()
	}
	return nil
}

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	FirstName   string `json:"first_name" validate:"omitempty,max=50"`
	LastName    string `json:"last_name" validate:"omitempty,max=50"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
	Password    string `json:"password" validate:"required,max=100"`
	RePassword  string `json:"re_password" validate:"required,max=100,eqfield=Password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required,min=3,max=50"`
	Password        string `json:"password" validate:"required,max=100"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" validate:"omitempty,max=50"`
	LastName    string `json:"last_name" validate:"omitempty,max=50"`
	Email       string `json:"email" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,max=100"`
	NewPassword     string `json:"new_password" validate:"required,max=100"`
	ReNewPassword   string `json:"re_new_password" validate:"required,max=100,eqfield=NewPassword"`
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
