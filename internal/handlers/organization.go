package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/repository"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/sirupsen/logrus"
)

// Organization management

type OrganizationHandler struct {
	repo *repository.OrganizationRepository
	log  *logrus.Logger
	svc  *services.OrganizationService
}

func NewOrganizationHandler(repo *repository.OrganizationRepository, log *logrus.Logger, svc *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		repo: repo,
		log:  log,
		svc:  svc,
	}
}

// Create creates a new organization
func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var org models.Organization
	if err := c.BodyParser(&org); err != nil {
		h.log.WithError(err).Error("failed to bind organization data")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	created, err := h.svc.Create(&org)
	if err != nil {
		h.log.WithError(err).Error("failed to create organization")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create organization"})
	}

	return c.Status(http.StatusCreated).JSON(created)
}

// Get retrieves an organization by ID
func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.WithError(err).Error("invalid organization ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid organization ID"})
	}
	org, err := h.svc.GetByID(uint(id))
	if err != nil {
		h.log.WithError(err).Error("failed to get organization")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "organization not found"})
	}

	return c.Status(http.StatusOK).JSON(org)
}

// List retrieves all organizations with optional filtering
func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	orgs, err := h.svc.GetAll()
	if err != nil {
		h.log.WithError(err).Error("failed to list organizations")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve organizations"})
	}

	return c.Status(http.StatusOK).JSON(orgs)
}

// Update updates an organization
func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var org models.Organization
	if err := c.BodyParser(&org); err != nil {
		h.log.WithError(err).Error("failed to bind organization data")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}

	org.ID = id
	updated, err := h.svc.Update(&org)
	if err != nil {
		h.log.WithError(err).Error("failed to update organization")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update organization"})
	}

	return c.Status(http.StatusOK).JSON(updated)
}

// Delete removes an organization
func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "organization ID is required"})
	}

	if err := h.svc.Delete(id); err != nil {
		h.log.WithError(err).Error("failed to delete organization")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete organization"})
	}

	return c.Status(http.StatusNoContent).SendString("")
}
