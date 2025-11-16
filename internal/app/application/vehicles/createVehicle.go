package vehicles

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type CreateVehicleInput struct {
	Plate       string
	Brand       string
	Model       string
	Color       string
	VehicleType string
}

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
)

type CreateVehicleUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewCreateVehicleUseCase(
	vehicleRepo repository.VehicleRepository,
) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{
		vehicleRepo: vehicleRepo,
	}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, input CreateVehicleInput) error {
	if err := validateVehicleInput(input); err != nil {
		return err
	}

	vehicle := entities.Vehicle{
		Plate:       input.Plate,
		Brand:       input.Brand,
		Model:       input.Model,
		Color:       input.Color,
		VehicleType: input.VehicleType,
	}

	if err := uc.vehicleRepo.CreateVehicleTx(&vehicle); err != nil {
		return err
	}

	return nil
}

func validateVehicleInput(input CreateVehicleInput) error {
	if input.Plate == "" || input.Brand == "" || input.Model == "" || input.Color == "" || input.VehicleType == "" {
		return ErrInvalidInput
	}
	return nil
}
