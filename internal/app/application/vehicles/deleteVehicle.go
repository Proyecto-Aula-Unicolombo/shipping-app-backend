package vehicles

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrVehicleNotFoundDelete = errors.New("Vehículo no encontrado")
)

type DeleteVehicleUseCase struct {
	repo repository.VehicleRepository
}

func NewDeleteVehicleUseCase(repo repository.VehicleRepository) *DeleteVehicleUseCase {
	return &DeleteVehicleUseCase{repo: repo}
}

func (uc *DeleteVehicleUseCase) Execute(id uint) error {
	// 1. Validar ID
	if id == 0 {
		return ErrInvalidID
	}

	// 2. Verificar que el vehículo existe
	_, err := uc.repo.GetVehicleByID(id)
	if err != nil {
		return ErrVehicleNotFound
	}

	// 3. Eliminar vehículo
	// TODO: Agregar validación - Solo COORDINADOR puede eliminar
	return uc.repo.DeleteVehicle(id)
}