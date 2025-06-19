package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/services"
)

type OCPPHandler struct {
	svc *services.ChargePointService
	log *logrus.Logger
}

func GocsmsOCPPHandler(svc *services.ChargePointService, log *logrus.Logger) *OCPPHandler {
	return &OCPPHandler{svc: svc, log: log}
}

func (h *OCPPHandler) HandleMessage(ctx context.Context, chargePointID string, msg []byte) ([]byte, error) {
	var ocppMsg OCPPMessage
	if err := json.Unmarshal(msg, &ocppMsg); err != nil {
		h.log.Error("Invalid OCPP message: ", err)
		return h.createErrorResponse(ocppMsg.UniqueID, "FormationViolation", "Invalid JSON")
	}

	if ocppMsg.MessageTypeID != Call {
		return h.createErrorResponse(ocppMsg.UniqueID, "NotSupported", "Only CALL messages supported")
	}

	switch ocppMsg.Action {
	case "BootNotification":
		return h.handleBootNotification(ctx, chargePointID, ocppMsg)
	case "Heartbeat":
		return h.handleHeartbeat(ctx, chargePointID, ocppMsg)
	case "StatusNotification":
		return h.handleStatusNotification(ctx, chargePointID, ocppMsg)
	default:
		return h.createErrorResponse(ocppMsg.UniqueID, "NotSupported", fmt.Sprintf("Action %s not supported", ocppMsg.Action))
	}
}

func (h *OCPPHandler) handleBootNotification(ctx context.Context, chargePointID string, msg OCPPMessage) ([]byte, error) {
	var req BootNotificationRequest
	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		return h.createErrorResponse(msg.UniqueID, "FormationViolation", "Invalid payload")
	}

	h.log.Infof("Received BootNotification from %s: %+v", chargePointID, req)
	err := h.svc.Register(ctx, &models.ChargePoint{
		ID:            chargePointID,
		SerialNumber:  req.ChargePointSerialNumber,
		Status:        "Available",
		LastHeartbeat: time.Now(),
	})
	if err != nil {
		h.log.Error("Failed to register charge point: ", err)
		return h.createErrorResponse(msg.UniqueID, "InternalError", err.Error())
	}

	resp := BootNotificationResponse{
		Status:      "Accepted",
		CurrentTime: time.Now(),
		Interval:    60,
	}
	return h.createResponse(msg.UniqueID, resp)
}

func (h *OCPPHandler) handleHeartbeat(ctx context.Context, chargePointID string, msg OCPPMessage) ([]byte, error) {
	h.log.Infof("Received Heartbeat from %s", chargePointID)
	err := h.svc.UpdateStatus(ctx, chargePointID, "Available")
	if err != nil {
		h.log.Error("Failed to update heartbeat: ", err)
		return h.createErrorResponse(msg.UniqueID, "InternalError", err.Error())
	}

	resp := HeartbeatResponse{
		CurrentTime: time.Now(),
	}
	return h.createResponse(msg.UniqueID, resp)
}

func (h *OCPPHandler) handleStatusNotification(ctx context.Context, chargePointID string, msg OCPPMessage) ([]byte, error) {
	var req StatusNotificationRequest
	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		return h.createErrorResponse(msg.UniqueID, "FormationViolation", "Invalid payload")
	}

	h.log.Infof("Received StatusNotification from %s: %+v", chargePointID, req)
	err := h.svc.UpdateStatus(ctx, chargePointID, req.Status)
	if err != nil {
		h.log.Error("Failed to update status: ", err)
		return h.createErrorResponse(msg.UniqueID, "InternalError", err.Error())
	}

	resp := StatusNotificationResponse{}
	return h.createResponse(msg.UniqueID, resp)
}

func (h *OCPPHandler) createResponse(uniqueID string, payload interface{}) ([]byte, error) {
	resp := OCPPMessage{
		MessageTypeID: CallResult,
		UniqueID:      uniqueID,
		Payload:       json.RawMessage{},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp.Payload = data
	return json.Marshal(resp)
}

func (h *OCPPHandler) createErrorResponse(uniqueID, errorCode, errorMessage string) ([]byte, error) {
	resp := OCPPMessage{
		MessageTypeID: CallError,
		UniqueID:      uniqueID,
		ErrorCode:     errorCode,
		ErrorMessage:  errorMessage,
	}
	return json.Marshal(resp)
}
