package entities

type Vehicle struct {
	ID          uint
	Plate       string
	Brand       string
	Model       string
	Color       string
	VehicleType string

	AssignedDriverName     string
	AssignedDriverLastName string
}
