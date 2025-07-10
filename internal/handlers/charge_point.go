package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/mutoulbj/gocsms/internal/dto"
	"github.com/mutoulbj/gocsms/internal/enums"
	"github.com/mutoulbj/gocsms/internal/middleware"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/mutoulbj/gocsms/internal/utils"
)

// @title Charge Point API
// @version 1.0
// @description API for managing charge points
// @BasePath /api/v1

type ChargePointHandler struct {
	svc     *services.ChargePointService
	authSvc *services.AuthService
	redis   *redis.Client
	log     *logrus.Logger
}

func NewChargePointHandler(svc *services.ChargePointService, authSvc *services.AuthService, redis *redis.Client, log *logrus.Logger) *ChargePointHandler {
	return &ChargePointHandler{svc: svc, authSvc: authSvc, redis: redis, log: log}
}

func (h *ChargePointHandler) RegisterRoutes(app fiber.Router) {
	cp := app.Group("/chargepoints", middleware.Auth(h.authSvc, h.redis, h.log))

	cp.Post("/", h.Create)                // @Summary Register a new charge point
	cp.Get("/:id", h.GetByID)             // @Summary Get charge point by ID
	cp.Put("/:id/status", h.UpdateStatus) // @Summary Update charge point status
}

// @Summary Create(Register) a new charge point
// @Description Register a new charge point
// @Tags ChargePoints
// @Accept json
// @Produce json
// @Param chargepoint body models.ChargePoint true "Charge point data"
// @Success 201 {object} models.ChargePoint
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /chargepoints [post]
func (h *ChargePointHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateChargePointRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	status := enums.ChargePointStatusUnknown
	if req.Status != "" {
		status = enums.ChargePointStatus(req.Status)
		if !status.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
		}
	}

	cp := models.ChargePoint{
		Name:         req.Name,
		Code:         req.Code,
		Status:       status,
		SerialNumber: utils.GenerateSerialNumber(),
	}

	if err := h.svc.Register(c.Context(), &cp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(cp)
}

// @Summary Get charge point by ID
// @Description Retrieve a charge point details by ID
// @Tags ChargePoints
// @Accept json
// @Produce json
// @Param id path string true "Charge point ID"
// @Success 200 {object} models.ChargePoint
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /chargepoints/{id} [get]
func (h *ChargePointHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	cp, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Charge point not found"})
	}
	return c.JSON(cp)
}

// @Summary Update charge point status
// @Description Update the status of a charge point
// @Tags ChargePoints
// @Accept json
// @Produce json
// @Param id path string true "Charge Point ID"
// @Param status body map[string]string true "Status"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /chargepoints/{id}/status [put]
func (h *ChargePointHandler) UpdateStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	uuidID, err := utils.ParseUUID(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID"})
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.svc.UpdateStatus(c.Context(), uuidID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Status updated"})
}
