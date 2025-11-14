package vehicles

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

type UpdateVehicleUseCase struct {
	repo repository.VehicleRepository
}

func NewUpdateVehicleUseCase(repo repository.VehicleRepository) *UpdateVehicleUseCase {
	return &UpdateVehicleUseCase{repo: repo}
}

type UpdateVehicleInput struct {
	ID          uint
	Plate       string
	Brand       string
	Model       string
	Color       string
	VehicleType string
}

func (uc *UpdateVehicleUseCase) Execute(input UpdateVehicleInput) error {
	if input.ID == 0 {
		return errors.New("ID inválido")
	}

	existingVehicle, err := uc.repo.GetVehicleByID(input.ID)
	if err != nil {
		return errors.New("vehículo no encontrado")
	}

	if input.Plate != "" {
		existingVehicle.Plate = input.Plate
	}
	if input.Brand != "" {
		existingVehicle.Brand = input.Brand
	}
	if input.Model != "" {
		existingVehicle.Model = input.Model
	}
	if input.Color != "" {
		existingVehicle.Color = input.Color
	}
	if input.VehicleType != "" {
		existingVehicle.VehicleType = input.VehicleType
	}

	
	return uc.repo.UpdateVehicle(existingVehicle)
}