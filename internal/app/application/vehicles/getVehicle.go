package vehicles

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrVehicleNotFound = errors.New("Vehículo no registrado")
	ErrInvalidID       = errors.New("ID inválido")
)

type VehiclesOutput struct {
	ID          uint
	Plate       string
	Brand       string
	Model       string
	Color       string
	VehicleType string
}

type GetVehicle struct {
	repo repository.VehicleRepository
}

func NewGetVehicle(repo repository.VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (uc *GetVehicle) Execute(ctx context.Context, id uint) (*VehiclesOutput, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrVehicleNotFound
	}

	vehiclesOutput := &VehiclesOutput{
		ID:          vehicle.ID,
		Plate:       vehicle.Plate,
		Brand:       vehicle.Brand,
		Model:       vehicle.Model,
		Color:       vehicle.Color,
		VehicleType: vehicle.VehicleType,
	}

	return vehiclesOutput, nil
}
