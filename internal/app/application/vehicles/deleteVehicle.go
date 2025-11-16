package vehicles

import (
	"context"
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

func (uc *DeleteVehicleUseCase) Execute(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrVehicleNotFound
	}

	return uc.repo.DeleteVehicle(id)
}
