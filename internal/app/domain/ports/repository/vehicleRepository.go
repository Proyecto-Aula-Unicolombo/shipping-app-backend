package repository

import (
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type VehicleRepository interface {
	CreateVehicleTx(tx *sql.Tx, vehicle *entities.Vehicle) error
	GetVehicleByID(id uint) (*entities.Vehicle, error)
	DeleteVehicle(id uint) error
	UpdateVehicle(vehicle *entities.Vehicle) error
	GetAllVehicles() ([]*entities.Vehicle, error)
	ListVehicles(limit, offset int, PlateOrBrand string) ([]*entities.Vehicle, error)
}
