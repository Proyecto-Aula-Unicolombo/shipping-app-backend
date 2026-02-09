package vehicles

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrVehicleNotFoundDelete  = errors.New("Vehículo no encontrado")
	ErrVehicleHasActiveOrders = errors.New("vehicle has active orders and cannot be deleted")
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

	hasActiveOrders, err := uc.repo.HasActiveVehicleInOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking active orders: %w", err)
	}

	if hasActiveOrders {
		return ErrVehicleHasActiveOrders
	}
	return uc.repo.DeleteVehicle(id)
}
