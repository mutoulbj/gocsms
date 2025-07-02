package models

import "time"

type APIResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorData  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Success   bool        `json:"success"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Pagination struct {
	Items      any   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *User     `json:"user"`
}
