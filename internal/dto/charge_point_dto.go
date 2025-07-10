package dto

type CreateChargePointRequest struct {
	Name               string `json:"name" validate:"required"`
	Code               string `json:"code" validate:"required"`
	Status             string `json:"status" validate:"required,oneof=available unavailable charging paused reserved"`
	SerdialNumber      string `json:"serial_number" validate:"required"`
	OcppVersion        string `json:"ocpp_version" validate:"required"`
	Model              string `json:"model" validate:"required"`
	Vendor             string `json:"vendor" validate:"required"`
	ChargeStationId    string `json:"charge_station_id" validate:"required,uuid"`
	RegistrationStatus string `json:"registration_status" validate:"required,oneof=accepted rejected pending"`
}
