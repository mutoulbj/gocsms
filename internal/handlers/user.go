package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mutoulbj/gocsms/internal/dto"
	"github.com/mutoulbj/gocsms/internal/middleware"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/mutoulbj/gocsms/internal/utils"
	"github.com/mutoulbj/gocsms/pkg/response"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	svc      *services.UserService
	authSvc  *services.AuthService
	log      *logrus.Logger
	res      response.APIResponseInterface
	redis    *redis.Client
	validate *validator.Validate
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	svc *services.UserService,
	authSvc *services.AuthService,
	log *logrus.Logger,
	res response.APIResponseInterface,
	redis *redis.Client,
) *UserHandler {
	return &UserHandler{
		svc:      svc,
		authSvc:  authSvc,
		log:      log,
		res:      res,
		redis:    redis,
		validate: validator.New(),
	}
}

// RegisterRoutes registers the user-related routes with the provided router
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	user := router.Group("/users")
	// Register user routes
	user.Post("/register", h.CreateUser)                                        // Create a new user
	user.Get("/:id", middleware.Auth(h.authSvc, h.redis, h.log), h.GetUserById) // Get user by ID
	user.Get("", middleware.Auth(h.authSvc, h.redis, h.log), h.ListUsers)       // List all users
	user.Put("/:id", middleware.Auth(h.authSvc, h.redis, h.log), h.UpdateUser)  // Update user by ID
}

// GetUserById retrieves a user by their ID
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.res.Error(c, http.StatusBadRequest, "invalid user ID", "params error", "user ID is required")
	}
	// Parse the ID to a UUID
	uid, err := utils.ParseUUID(id)
	if err != nil {
		h.log.WithError(err).Error("Invalid user ID format")
		return h.res.Error(c, http.StatusBadRequest, "invalid user ID", "params error", "user ID is required")
	}

	user, err := h.svc.GetUserById(c.Context(), uid)
	if err != nil {
		h.log.WithError(err).Error("Failed to get user by ID")
		return h.res.Error(c, http.StatusInternalServerError, "failed to get user", "internal error", err.Error())
	}

	return h.res.Success(c, "User retrieved successfully", user)
}

// ListUsers lists all users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// search by username or email
	username := c.Query("username", "")
	email := c.Query("email", "")

	// page and pageSize can be added later for pagination
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	users, err := h.svc.ListUsers(c.Context(), username, email, page, pageSize)
	if err != nil {
		h.log.WithError(err).Error("Failed to list users")
		return h.res.ErrorHandler(c, err)
	}

	return h.res.Success(c, "Users listed successfully", users)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to parse user data")
		return h.res.Error(c, http.StatusBadRequest, "invalid user data", "params error", err.Error())
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		h.log.WithError(err).Error("Validation failed for user data")
		return h.res.Error(c, http.StatusBadRequest, "invalid user data", "params error", err.Error())
	}

	createdUser, err := h.svc.CreateUser(c.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to create user")
		return h.res.ErrorHandler(c, err)
	}

	return h.res.Created(c, "User created successfully", createdUser)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.res.Error(c, http.StatusBadRequest, "invalid user ID", "params error", "user ID is required")
	}

	// Parse the ID to a UUID
	uid, err := utils.ParseUUID(id)
	if err != nil {
		h.log.WithError(err).Error("Invalid user ID format")
		return h.res.Error(c, http.StatusBadRequest, "invalid user ID", "params error", "user ID is required")
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to parse update request")
		return h.res.Error(c, http.StatusBadRequest, "invalid update data", "params error", err.Error())
	}

	updatedUser, err := h.svc.UpdateUser(c.Context(), uid, &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to update user")
		return h.res.ErrorHandler(c, err)
	}

	return h.res.Success(c, "User updated successfully", updatedUser)
}
