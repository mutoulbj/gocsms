package enums

type ChargePointStatus string

const (
	ChargePointStatusUnknown       ChargePointStatus = "UNKNOWN"
	ChargePointStatusAvailable     ChargePointStatus = "AVAILABLE"
	ChargePointStatusPreparing     ChargePointStatus = "PREPARING"
	ChargePointStatusCharging      ChargePointStatus = "CHARGING"
	ChargePointStatusFinishing     ChargePointStatus = "FINISHING"
	ChargePointStatusSuspendedEVSE ChargePointStatus = "SUSPENDED_EVSE"
	ChargePointStatusSuspendedEV   ChargePointStatus = "SUSPENDED_EV"
	ChargePointStatusUnavailable   ChargePointStatus = "UNAVAILABLE"
	ChargePointStatusFaulted       ChargePointStatus = "FAULTED"
	ChargePointStatusOffline       ChargePointStatus = "OFFLINE"
)

func (s ChargePointStatus) IsValid() bool {
	switch s {
	case ChargePointStatusUnknown, ChargePointStatusAvailable, ChargePointStatusPreparing, ChargePointStatusCharging,
		ChargePointStatusFinishing, ChargePointStatusSuspendedEVSE, ChargePointStatusSuspendedEV,
		ChargePointStatusUnavailable, ChargePointStatusFaulted, ChargePointStatusOffline:
		return true
	default:
		return false
	}
}
