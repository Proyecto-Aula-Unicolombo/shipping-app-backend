package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type CreateOrderInput struct {
	Observation *string
	DriverID    uint
	VehicleID   uint
	PackageIDs  []uint
}

type CreateOrderOutput struct {
	ID       uint
	Status   string
	CreateAt time.Time
}

var (
	ErrInvalidOrderInput   = errors.New("invalid order input")
	ErrDriverNotFound      = errors.New("driver not found")
	ErrVehicleNotFound     = errors.New("vehicle not found")
	ErrNoPackagesProvided  = errors.New("no packages provided for order")
	ErrPackageNotAvailable = errors.New("one or more packages are not available")
)

type CreateOrderUseCase struct {
	orderRepo   repository.OrderRepository
	driverRepo  repository.DriverRepository
	vehicleRepo repository.VehicleRepository
	packageRepo repository.PackageRepository
	txProvider  repository.TxProvider
}

func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository,
	packageRepo repository.PackageRepository,
	txProvider repository.TxProvider,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:   orderRepo,
		driverRepo:  driverRepo,
		vehicleRepo: vehicleRepo,
		packageRepo: packageRepo,
		txProvider:  txProvider,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// Validar input
	if err := validateCreateOrderInput(input); err != nil {
		return nil, err
	}

	// Verificar que el conductor existe
	driver, err := uc.driverRepo.GetByID(ctx, input.DriverID)
	if err != nil || driver == nil {
		return nil, ErrDriverNotFound
	}

	// Verificar que el vehículo existe
	vehicle, err := uc.vehicleRepo.GetByID(ctx, input.VehicleID)
	if err != nil || vehicle == nil {
		return nil, ErrVehicleNotFound
	}

	// Verificar que todos los paquetes existen y están disponibles (sin orden asignada)
	if err := uc.validatePackagesAvailability(ctx, input.PackageIDs); err != nil {
		return nil, err
	}

	// Iniciar transacción
	tx, err := uc.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = uc.txProvider.RollbackTx(ctx, tx)
		}
	}()

	// Crear orden
	now := time.Now()
	order := &entities.Order{
		CreateAt:    now,
		AssignedAt:  &now,
		Observation: input.Observation,
		Status:      "Pendiente",
		DriverID:    input.DriverID,
		VehicleID:   input.VehicleID,
	}

	if err := uc.orderRepo.Create(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// Asignar paquetes a la orden
	if err := uc.assignPackagesToOrder(ctx, tx, order.ID, input.PackageIDs); err != nil {
		return nil, fmt.Errorf("assign packages to order: %w", err)
	}

	for _, pkgID := range input.PackageIDs {
		if err := uc.packageRepo.UpdatePackageStatusDelivery(ctx, tx, "asignado", pkgID); err != nil {
			return nil, fmt.Errorf("update package status delivery: %w", err)
		}
	}

	// Commit
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return &CreateOrderOutput{
		ID:       order.ID,
		Status:   order.Status,
		CreateAt: order.CreateAt,
	}, nil
}

func validateCreateOrderInput(input CreateOrderInput) error {
	if input.DriverID == 0 {
		return ErrInvalidOrderInput
	}
	if input.VehicleID == 0 {
		return ErrInvalidOrderInput
	}
	if len(input.PackageIDs) == 0 {
		return ErrNoPackagesProvided
	}
	return nil
}

func (uc *CreateOrderUseCase) validatePackagesAvailability(ctx context.Context, packageIDs []uint) error {
	for _, pkgID := range packageIDs {
		pkg, err := uc.packageRepo.GetByID(ctx, pkgID)
		if err != nil {
			return fmt.Errorf("package %d not found: %w", pkgID, err)
		}
		// Verificar que el paquete no tenga orden asignada
		if pkg.OrderID != nil && *pkg.OrderID != 0 {
			return ErrPackageNotAvailable
		}
	}
	return nil
}

func (uc *CreateOrderUseCase) assignPackagesToOrder(ctx context.Context, tx *sql.Tx, orderID uint, packageIDs []uint) error {
	for _, pkgID := range packageIDs {
		pkg, err := uc.packageRepo.GetByID(ctx, pkgID)
		if err != nil {
			return fmt.Errorf("get package %d: %w", pkgID, err)
		}

		// Asignar OrderID al paquete
		pkg.OrderID = &orderID

		// Aquí necesitarías un método Update en PackageRepository
		// Por ahora, ejecutamos un UPDATE directo
		query := `UPDATE packages SET idorder = $1 WHERE id = $2`
		if _, err := tx.ExecContext(ctx, query, orderID, pkgID); err != nil {
			return fmt.Errorf("update package %d with order: %w", pkgID, err)
		}
	}
	return nil
}
