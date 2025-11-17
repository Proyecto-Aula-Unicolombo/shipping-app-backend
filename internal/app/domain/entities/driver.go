package entities

type Driver struct {
	ID          uint
	PhoneNumber string
	LicenseNo   string
	IsActive    bool
	UserID      uint
	User        *User

	NumOrder uint
}
