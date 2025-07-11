package dto

import (
	"time"

	"github.com/mutoulbj/gocsms/internal/models"
)

type UserResponse struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	LastLoginAt    string `json:"last_login_at"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
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

func ToUserResponse(u *models.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:             u.ID.String(),
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		Status:         u.Status.String(),
		CreatedAt:      u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      u.UpdatedAt.Format(time.RFC3339),
		LastLoginAt:    u.LastLoginAt.Format(time.RFC3339),
	}
}

func ToUserListResponse(users []*models.User, total int64, page, pageSize int) *UserListResponse {
	if users == nil {
		users = []*models.User{}
	}
	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return &UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)), // Calculate total pages
	}
}
