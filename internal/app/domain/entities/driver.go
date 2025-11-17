package entities

type Driver struct {
	ID          uint
	PhoneNumber string
	LicenseNo   string
	UserID      uint
	User        *User

	NumOrder uint
}
