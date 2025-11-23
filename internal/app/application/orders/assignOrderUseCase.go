package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type AssignOrderInput struct {
	OrderID   uint
	DriverID  uint
	VehicleID uint
}

var (
	ErrOrderCannotBeReassigned = errors.New("order with status 'entregado' cannot be reassigned")
	ErrOrderNotFound           = errors.New("order not found")
)

type AssignOrderUseCase struct {
	orderRepo   repository.OrderRepository
	driverRepo  repository.DriverRepository
	vehicleRepo repository.VehicleRepository
	txProvider  repository.TxProvider
}

func NewAssignOrderUseCase(
	orderRepo repository.OrderRepository,
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository,
	txProvider repository.TxProvider,
) *AssignOrderUseCase {
	return &AssignOrderUseCase{
		orderRepo:   orderRepo,
		driverRepo:  driverRepo,
		vehicleRepo: vehicleRepo,
		txProvider:  txProvider,
	}
}

func (uc *AssignOrderUseCase) Execute(ctx context.Context, input AssignOrderInput) error {
	order, err := uc.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if err := uc.validateOrderCanBeAssigned(order); err != nil {
		return err
	}

	if err := uc.validateDriverExists(ctx, input.DriverID); err != nil {
		return err
	}

	if err := uc.validateVehicleExists(ctx, input.VehicleID); err != nil {
		return err
	}

	tx, err := uc.txProvider.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = uc.txProvider.RollbackTx(ctx, tx)
		}
	}()

	if order.DriverID != nil && *order.DriverID != 0 {
		if *order.DriverID != input.DriverID {
			if err := uc.releaseCurrentDriver(*order.DriverID); err != nil {
				return fmt.Errorf("release current driver: %w", err)
			}
		}
	}

	if err := uc.assignNewDriver(input.DriverID); err != nil {
		return fmt.Errorf("assign new driver: %w", err)
	}
	if err := uc.updateOrderAssignment(ctx, tx, order, input); err != nil {
		return fmt.Errorf("update order assignment: %w", err)
	}

	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

// validateOrderCanBeAssigned verifica que la orden puede ser asignada o reasignada.
// Solo las órdenes con estado "entregado" no pueden ser reasignadas.
func (uc *AssignOrderUseCase) validateOrderCanBeAssigned(order *entities.Order) error {
	if order.Status == "entregado" {
		return ErrOrderCannotBeReassigned
	}
	return nil
}

// validateDriverExists verifica que el conductor existe en el sistema.
func (uc *AssignOrderUseCase) validateDriverExists(ctx context.Context, driverID uint) error {
	driver, err := uc.driverRepo.GetByID(ctx, driverID)
	if err != nil {
		return fmt.Errorf("get driver: %w", err)
	}
	if driver == nil {
		return ErrDriverNotFound
	}
	return nil
}

// validateVehicleExists verifica que el vehículo existe en el sistema.
func (uc *AssignOrderUseCase) validateVehicleExists(ctx context.Context, vehicleID uint) error {
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}
	return nil
}

// releaseCurrentDriver marca al conductor actual como disponible (is_busy = false).
// Esto se llama cuando se reasigna una orden a otro conductor.
func (uc *AssignOrderUseCase) releaseCurrentDriver(currentDriverID uint) error {
	if err := uc.driverRepo.UpdateDriverStatus(currentDriverID, false); err != nil {
		return fmt.Errorf("update current driver status to available: %w", err)
	}
	return nil
}

// assignNewDriver marca al nuevo conductor como ocupado (is_busy = true).
func (uc *AssignOrderUseCase) assignNewDriver(newDriverID uint) error {
	if err := uc.driverRepo.UpdateDriverStatus(newDriverID, true); err != nil {
		return fmt.Errorf("update new driver status to busy: %w", err)
	}
	return nil
}

// updateOrderAssignment actualiza la orden con el nuevo conductor y vehículo.
// El método AssignDriverAndVehicle ya actualiza el estado a "asignada" automáticamente.
func (uc *AssignOrderUseCase) updateOrderAssignment(ctx context.Context, tx *sql.Tx, order *entities.Order, input AssignOrderInput) error {
	// Actualizar conductor y vehículo (esto también actualiza el estado a "asignada")
	if err := uc.orderRepo.AssignDriverAndVehicle(ctx, input.OrderID, input.DriverID, input.VehicleID); err != nil {
		return fmt.Errorf("assign driver and vehicle: %w", err)
	}

	return nil
}
