package drivers

type DriverDTO struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	LicenseNo   string `json:"num_licence"`
}

type DriverUpdateDTO struct {
	PhoneNumber string `json:"phone_number"`
	LicenseNo   string `json:"num_licence"`
}
