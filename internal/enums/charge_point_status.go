package enums

import (
	"strings"
)

type ChargePointStatus int

const (
    ChargePointStatusUnknown ChargePointStatus = iota
	ChargePointStatusAvailable
	ChargePointStatusUnavailable
	ChargePointStatusCharging
	ChargePointStatusPaused
	ChargePointStatusReserved
)

func (s ChargePointStatus) String() string {
	switch s {
	case ChargePointStatusAvailable:
		return "available"
	case ChargePointStatusUnavailable:
		return "unavailable"
	case ChargePointStatusCharging:
		return "charging"
	case ChargePointStatusPaused:
		return "paused"
	case ChargePointStatusReserved:
		return "reserved"
	default:
		return "unknown"
	}
}

func (s ChargePointStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *ChargePointStatus) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	switch str {
	case "available":
		*s = ChargePointStatusAvailable
	case "unavailable":
		*s = ChargePointStatusUnavailable
	case "charging":
		*s = ChargePointStatusCharging
	case "paused":
		*s = ChargePointStatusPaused
	case "reserved":
		*s = ChargePointStatusReserved
	case "":
		*s = ChargePointStatusUnknown // Handle empty string as unknown
	default:
		*s = ChargePointStatusUnknown
	}
	return nil
}

func ChargePointStatusFromString(status string) ChargePointStatus {
	switch strings.ToLower(status) {
	case "available":
		return ChargePointStatusAvailable
	case "unavailable":
		return ChargePointStatusUnavailable
	case "charging":
		return ChargePointStatusCharging
	case "paused":
		return ChargePointStatusPaused
	case "reserved":
		return ChargePointStatusReserved
	default:
		return ChargePointStatusUnknown
	}
}