package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/lib/pq"
)

type CreateOrderInput struct {
	Observation *string
	DriverID    *uint
	VehicleID   *uint
	PackageIDs  []uint
	TypeService string
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
	if input.DriverID != nil {
		driver, err := uc.driverRepo.GetByID(ctx, *input.DriverID)
		if err != nil || driver == nil {
			return nil, ErrDriverNotFound
		}
	}

	// Verificar que el vehículo existe
	if input.VehicleID != nil {
		vehicle, err := uc.vehicleRepo.GetByID(ctx, *input.VehicleID)
		if err != nil || vehicle == nil {
			return nil, ErrVehicleNotFound
		}
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
	var assignedAt *time.Time
	var status string

	if input.DriverID != nil && input.VehicleID != nil {
		// Orden completa - asignada
		status = "asignada"
		assignedAt = &now
	} else {
		// Orden sin asignar
		status = "pendiente"
		assignedAt = nil
	}
	order := &entities.Order{
		CreateAt:    now,
		AssignedAt:  assignedAt,
		Observation: input.Observation,
		Status:      status,
		TypeService: input.TypeService,
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

	packageStatus := "asignado"

	for _, pkgID := range input.PackageIDs {
		if err := uc.packageRepo.UpdatePackageStatusDelivery(ctx, tx, packageStatus, pkgID); err != nil {
			return nil, fmt.Errorf("update package status delivery: %w", err)
		}
	}

	if input.DriverID != nil {
		if err := uc.driverRepo.UpdateDriverStatus(*input.DriverID, true); err != nil {
			return nil, fmt.Errorf("update driver status: %w", err)
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
	if len(input.PackageIDs) == 0 {
		return ErrNoPackagesProvided
	}
	if input.TypeService == "" {
		return ErrInvalidOrderInput
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
	if len(packageIDs) == 0 {
		return nil
	}

	query := `UPDATE packages SET idorder = $1 WHERE id = ANY($2)`

	if _, err := tx.ExecContext(ctx, query, orderID, pq.Array(packageIDs)); err != nil {
		return fmt.Errorf("update packages with order: %w", err)
	}

	return nil
}
