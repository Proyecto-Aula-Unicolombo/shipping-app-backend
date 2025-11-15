package vehicles

import (
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrVehicleNotFound = errors.New("Vehículo no registrado")
	ErrInvalidID       = errors.New("ID inválido")
)

type GetVehicle struct {
	repo repository.VehicleRepository
}

func NewGetVehicle(repo repository.VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (uc *GetVehicle) Execute(id uint) (*entities.Vehicle, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	vehicle, err := uc.repo.GetVehicleByID(id)
	if err != nil {
		return nil, ErrVehicleNotFound
	}

	return vehicle, nil
}