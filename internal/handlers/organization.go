package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/mutoulbj/gocsms/pkg/response"
	"github.com/sirupsen/logrus"
)

// Organization management

type OrganizationHandler struct {
	log *logrus.Logger
	svc *services.OrganizationService
	res response.APIResponseInterface
}

func NewOrganizationHandler(log *logrus.Logger, svc *services.OrganizationService, res response.APIResponseInterface) *OrganizationHandler {
	return &OrganizationHandler{
		log: log,
		svc: svc,
		res: res,
	}
}

// Create creates a new organization
func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var req models.OrganizationCreateRequest

	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("failed to bind organization data")
		h.res.ErrorHandler(c, err)
		return nil
	}
	created, err := h.svc.CreateFromRequest(c.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("failed to create organization")
		h.res.ErrorHandler(c, err)
		return nil
	}

	return h.res.Created(c, "Organization created", created)
}

// Get retrieves an organization by ID
func (h *OrganizationHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.WithError(err).Error("invalid organization ID")
		h.res.Error(c, http.StatusBadRequest, "invalid organization ID", "params error", err.Error())
		return nil
	}
	org, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		h.log.WithError(err).Error("failed to get organization")
		h.res.Error(c, http.StatusNotFound, "organization not found", "not found", err.Error())
		return nil
	}
	return h.res.Success(c, "Organization retrieved", org)
}

// List retrieves all organizations with optional filtering
func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.log.WithError(err).Error("invalid page number")
		h.res.Error(c, http.StatusBadRequest, "invalid page number", "params error", err.Error())
		return nil
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		h.log.WithError(err).Error("invalid page size")
		h.res.Error(c, http.StatusBadRequest, "invalid page size", "params error", err.Error())
		return nil
	}

	orgs, total, err := h.svc.GetAll(c.Context(), page, pageSize)
	if err != nil {
		h.log.WithError(err).Error("failed to list organizations")
		h.res.Error(c, http.StatusInternalServerError, "failed to retrieve organizations", "internal error", err.Error())
		return nil
	}

	return h.res.Success(c, "Organizations retrieved", fiber.Map{
		"data":  orgs,
		"total": total,
		"page":  page,
	})
}

func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.WithError(err).Error("invalid organization ID")
		h.res.Error(c, http.StatusBadRequest, "invalid organization ID", "params error", err.Error())
		return nil
	}

	var org models.Organization
	if err := c.BodyParser(&org); err != nil {
		h.log.WithError(err).Error("failed to bind organization data")
		h.res.Error(c, http.StatusBadRequest, "invalid request payload", "params error", err.Error())
		return nil
	}
	updated, err := h.svc.Update(c.Context(), id, &org)
	if err != nil {
		h.log.WithError(err).Error("failed to update organization")
		h.res.Error(c, http.StatusInternalServerError, "failed to update organization", "internal error", err.Error())
		return nil
	}
	return h.res.Success(c, "Organization updated", updated)
}

func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.WithError(err).Error("invalid organization ID")
		return h.res.ErrorHandler(c, err)
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		h.log.WithError(err).Error("failed to delete organization")
		h.res.ErrorHandler(c, err)
		return nil
	}
	return h.res.Success(c, "Organization deleted", nil)
}
