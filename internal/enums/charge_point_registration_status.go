package enums

type ChargePointRegistrationStatus string

const (
	ChargePointRegistrationStatusUnknown  ChargePointRegistrationStatus = "UNKNOWN"
	ChargePointRegistrationStatusAccepted ChargePointRegistrationStatus = "ACCEPTED"
	ChargePointRegistrationStatusRejected ChargePointRegistrationStatus = "REJECTED"
	ChargePointRegistrationStatusPending  ChargePointRegistrationStatus = "PENDING"
)

func (s ChargePointRegistrationStatus) IsValid() bool {
	switch s {
	case ChargePointRegistrationStatusUnknown, ChargePointRegistrationStatusAccepted,
		ChargePointRegistrationStatusRejected, ChargePointRegistrationStatusPending:
		return true
	default:
		return false
	}
}
