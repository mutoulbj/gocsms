package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type APIResponseInterface interface {
	Success(c *fiber.Ctx, message string, data any) error
	Created(c *fiber.Ctx, message string, data any) error
	Error(c *fiber.Ctx, statusCode int, code, message string, details any) error
	ValidationError(c *fiber.Ctx, details any) error
	Unauthorized(c *fiber.Ctx, message string) error
	NotFound(c *fiber.Ctx, message string) error
	Paginated(c *fiber.Ctx, message string, items any, page, pageSize int, total int64) error
	ErrorHandler(c *fiber.Ctx, err error) error
}

type APIResponse struct {
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorData `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
	Success   bool       `json:"success"`
}

type APIResponseHandler struct {
	log *logrus.Logger
}

func NewAPIResponse(log *logrus.Logger) APIResponseInterface {
	return &APIResponseHandler{
		log: log,
	}
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
}

func (r *APIResponseHandler) Success(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (r *APIResponseHandler) Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (r *APIResponseHandler) Error(c *fiber.Ctx, statusCode int, code, message string, details any) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: "Request failed",
		Error: &ErrorData{
			Code:    statusCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	})
}

func (r *APIResponseHandler) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	// Log the error (optional)
	r.log.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		// If the error is a fiber.Error, use its status code
		code = e.Code
	}
	// Return a generic error response
	return r.Error(c, code, "INTERNAL_ERROR", err.Error(), nil)
}

func (r *APIResponseHandler) ValidationError(c *fiber.Ctx, details any) error {
	return r.Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

func (r *APIResponseHandler) Unauthorized(c *fiber.Ctx, message string) error {
	return r.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func (r *APIResponseHandler) NotFound(c *fiber.Ctx, message string) error {
	return r.Error(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

func (r *APIResponseHandler) Paginated(c *fiber.Ctx, message string, items any, page, pageSize int, total int64) error {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	data := Pagination{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	return r.Success(c, message, data)
}
