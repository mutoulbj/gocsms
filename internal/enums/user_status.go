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
