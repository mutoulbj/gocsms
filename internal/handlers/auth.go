package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mutoulbj/gocsms/internal/middleware"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authSvc *services.AuthService
	log     *logrus.Logger
}

func NewAuthHandler(authSvc *services.AuthService, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		log:     log,
	}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)                                             // @Summary User login
	auth.Post("/refresh", middleware.Auth(h.authSvc, nil, h.log), h.Refresh) // @Summary Refresh JWT token
	auth.Post("/logout", middleware.Auth(h.authSvc, nil, h.log), h.Logout)   // @Summary User logout
	auth.Post("/kick", middleware.Auth(h.authSvc, nil, h.log), h.Kick)       // @Summary Kick user session
}

// @Summary Login to get access and refresh tokens
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Credentials"
// @Success 200 {object} services.TokenPair
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tokens, err := h.authSvc.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	return c.JSON(tokens)
}

// @Summary Logout user
// @Description Invalidate user session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /auth/logout [post]
// @Security BearerAuth
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tokenID := c.Locals("token_id").(string)

	if err := h.authSvc.Logout(c.Context(), userID, tokenID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User logged out successfully"})
}

// @Summary Refresh access token
// @Description Use refresh token to get a new token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "Refresh Token"
// @Success 200 {object} services.TokenPair
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tokens, err := h.authSvc.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}
	return c.JSON(tokens)
}

// @Summary Kick a user session
// @Description Invalidate a user’s session by user_id and token_id
// @Tags Auth
// @Accept json
// @Produce json
// @Param session body map[string]string true "Session Info"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /auth/kick [post]
// @Security BearerAuth
func (h *AuthHandler) Kick(c *fiber.Ctx) error {
	var req struct {
		UserID  string `json:"user_id"`
		TokenID string `json:"token_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.authSvc.KickUser(c.Context(), req.UserID, req.TokenID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User session kicked"})
}
