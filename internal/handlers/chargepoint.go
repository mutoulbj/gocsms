package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/services"
	"github.com/mutoulbj/gocsms/internal/utils"
)

// @title Charge Point API
// @version 1.0
// @description API for managing charge points
// @BasePath /api/v1
type ChargePointHandler struct {
	svc *services.ChargePointService
	log *logrus.Logger
}

func GocsmsChargePointHandler(svc *services.ChargePointService, log *logrus.Logger) *ChargePointHandler {
	return &ChargePointHandler{svc: svc, log: log}
}

func (h *ChargePointHandler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	cp := v1.Group("/chargepoints")

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
	var req models.CreateChargePointRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	status := req.Status
	if status == "" {
		status = "Available"
	}
	cp := models.ChargePoint{
		Name:         req.Name,
		Code:         req.Code,
		Status:       status,                       // Default to "Available" if not provided
		SerialNumber: utils.GenerateSerialNumber(), // Generate or assign a serial number as needed
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.svc.UpdateStatus(c.Context(), id, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Status updated"})
}
