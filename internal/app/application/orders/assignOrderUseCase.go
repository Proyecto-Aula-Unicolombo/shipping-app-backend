package orders

import (
	"context"
	"errors"
	"fmt"

	"shipping-app/internal/app/domain/ports/repository"
)

type AssignOrderInput struct {
	OrderID   uint
	DriverID  uint
	VehicleID uint
}

var (
	ErrOrderAlreadyAssigned = errors.New("order already assigned")
)

type AssignOrderUseCase struct {
	orderRepo   repository.OrderRepository
	driverRepo  repository.DriverRepository
	vehicleRepo repository.VehicleRepository
}

func NewAssignOrderUseCase(
	orderRepo repository.OrderRepository,
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository,
) *AssignOrderUseCase {
	return &AssignOrderUseCase{
		orderRepo:   orderRepo,
		driverRepo:  driverRepo,
		vehicleRepo: vehicleRepo,
	}
}

func (uc *AssignOrderUseCase) Execute(ctx context.Context, input AssignOrderInput) error {
	// Verificar que la orden existe
	order, err := uc.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// Verificar que la orden no esté ya asignada
	if order.Status != "Pendiente" {
		return ErrOrderAlreadyAssigned
	}

	// Verificar que el conductor existe
	driver, err := uc.driverRepo.GetByID(ctx, input.DriverID)
	if err != nil || driver == nil {
		return ErrDriverNotFound
	}

	// Verificar que el vehículo existe
	vehicle, err := uc.vehicleRepo.GetByID(ctx, input.VehicleID)
	if err != nil || vehicle == nil {
		return ErrVehicleNotFound
	}

	// Asignar conductor y vehículo
	if err := uc.orderRepo.AssignDriverAndVehicle(ctx, input.OrderID, input.DriverID, input.VehicleID); err != nil {
		return fmt.Errorf("assign driver and vehicle: %w", err)
	}

	return nil
}
