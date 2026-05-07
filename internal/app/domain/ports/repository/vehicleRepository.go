package repository

import (
	"context"
	"shipping-app/internal/app/domain/entities"
)

type VehicleRepository interface {
	CreateVehicleTx(vehicle *entities.Vehicle) error
	GetByID(ctx context.Context, id uint) (*entities.Vehicle, error)
	HasActiveVehicleInOrder(ctx context.Context, id uint) (bool, error)
	DeleteVehicle(id uint) error
	UpdateVehicle(vehicle *entities.Vehicle) error
	ListVehicles(limit, offset int, PlateOrBrand, vehicleType string) ([]*entities.Vehicle, error)
	CountVehicles(PlateOrBrand, vehicleType string) (int64, error)
	ListVehiclesUnassigned() ([]*entities.Vehicle, error)
}
