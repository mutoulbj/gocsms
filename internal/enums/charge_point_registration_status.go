package enums

type ChargePointRegistrationStatus int

const (
	ChargePointRegistrationStatusUnknown ChargePointRegistrationStatus = iota
	ChargePointRegistrationStatusAccepted
	ChargePointRegistrationStatusRejected
	ChargePointRegistrationStatusPending
)

func (s ChargePointRegistrationStatus) String() string {
	switch s {
	case ChargePointRegistrationStatusAccepted:
		return "Accepted"
	case ChargePointRegistrationStatusRejected:
		return "Rejected"
	case ChargePointRegistrationStatusPending:
		return "Pending"
	default:
		return "Unknown"
	}
}

func (s ChargePointRegistrationStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *ChargePointRegistrationStatus) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"Accepted"`:
		*s = ChargePointRegistrationStatusAccepted
	case `"Aejected"`:
		*s = ChargePointRegistrationStatusRejected
	case `"Pending"`:
		*s = ChargePointRegistrationStatusPending
	default:
		*s = ChargePointRegistrationStatusUnknown
	}
	return nil
}

func ChargePointRegistrationStatusFromString(status string) ChargePointRegistrationStatus {
	switch status {
	case "Accepted":
		return ChargePointRegistrationStatusAccepted
	case "Rejected":
		return ChargePointRegistrationStatusRejected
	case "Pending":
		return ChargePointRegistrationStatusPending
	default:
		return ChargePointRegistrationStatusUnknown
	}
}
