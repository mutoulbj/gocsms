package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func Cache() fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
		Methods:      []string{fiber.MethodGet},
	})
}