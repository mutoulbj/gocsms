package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func Auth(authSvc *services.AuthService, redisClient *redis.Client, log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is required"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}

		log.Debug("Validating token: ", parts[1])
		claims, err := authSvc.ValidateToken(parts[1], false)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// check session in redis
		sessionKey := "session:" + claims.UserID + ":" + claims.TokenID
		if _, err := redisClient.Get(c.Context(), sessionKey).Result(); err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session not found"})
		} else if err != nil {
			log.Error("Failed to check session in Redis: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		// set user context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		return c.Next()
	}
}
