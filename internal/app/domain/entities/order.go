package entities

import "time"

type Order struct {
	ID          uint
	CreateAt    time.Time
	AssignedAt  *time.Time
	Observation *string
	Status      string
	DriverID    uint
	VehicleID   uint
	Driver      *Driver
	Vehicle     *Vehicle
}
