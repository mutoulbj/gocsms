package enums

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended, UserStatusBanned:
		return true
	default:
		return false
	}
}

func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "Active"
	case UserStatusSuspended:
		return "Suspended"
	case UserStatusBanned:
		return "Banned"
	default:
		return "Unknown"
	}
}
