package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/sirupsen/logrus"
)

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func Error(c *fiber.Ctx, statusCode int, code, message string, details any) error {
	return c.Status(statusCode).JSON(models.APIResponse{
		Success: false,
		Message: "Request failed",
		Error: &models.ErrorData{
			Code:    statusCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	// Log the error (optional)
	logrus.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		// If the error is a fiber.Error, use its status code
		code = e.Code
	}
	// Return a generic error response
	return Error(c, code, "INTERNAL_ERROR", err.Error(), nil)
}

func ValidationError(c *fiber.Ctx, details any) error {
	return Error(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

func Paginated(c *fiber.Ctx, message string, items any, page, pageSize int, total int64) error {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	data := models.Pagination{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	return Success(c, message, data)
}
