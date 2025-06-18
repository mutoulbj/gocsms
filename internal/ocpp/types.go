package ocpp

import (
	"encoding/json"
	"time"
)

// MessageType defines OCPP message types
type MessageType int

const (
	Call        MessageType = 2
	CallResult  MessageType = 3
	CallError   MessageType = 4
)

// OCPPMessage represents a generic OCPP message
type OCPPMessage struct {
	MessageTypeID MessageType     `json:"messageTypeId"`
	UniqueID      string          `json:"uniqueId"`
	Action        string          `json:"action,omitempty"`
	Payload       json.RawMessage `json:"payload"`
	ErrorCode     string          `json:"errorCode,omitempty"`
	ErrorMessage  string          `json:"errorMessage,omitempty"`
}

// BootNotificationRequest for OCPP 1.6
type BootNotificationRequest struct {
	ChargePointVendor       string `json:"chargePointVendor"`
	ChargePointModel        string `json:"chargePointModel"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty"`
}

// BootNotificationResponse for OCPP 1.6
type BootNotificationResponse struct {
	Status      string    `json:"status"` // Accepted, Rejected, Pending
	CurrentTime time.Time `json:"currentTime"`
	Interval    int       `json:"interval"`
}

// HeartbeatRequest for OCPP 1.6
type HeartbeatRequest struct {
	// Empty payload as per OCPP 1.6
}

// HeartbeatResponse for OCPP 1.6
type HeartbeatResponse struct {
	CurrentTime time.Time `json:"currentTime"`
}

// StatusNotificationRequest for OCPP 1.6
type StatusNotificationRequest struct {
	ConnectorID int       `json:"connectorId"`
	Status      string    `json:"status"` // Available, Preparing, Charging, etc.
	ErrorCode   string    `json:"errorCode"`
	Timestamp   time.Time `json:"timestamp"`
}

// StatusNotificationResponse for OCPP 1.6
type StatusNotificationResponse struct {
	// Empty payload as per OCPP 1.6
}