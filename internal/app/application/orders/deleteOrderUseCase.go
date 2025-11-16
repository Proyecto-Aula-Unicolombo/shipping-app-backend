package orders

import (
	"context"
	"errors"
	"fmt"

	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrCannotDeleteOrder = errors.New("cannot delete order with status other than Pendiente")
)

type DeleteOrderUseCase struct {
	orderRepo   repository.OrderRepository
	packageRepo repository.PackageRepository
	txProvider  repository.TxProvider
}

func NewDeleteOrderUseCase(
	orderRepo repository.OrderRepository,
	packageRepo repository.PackageRepository,
	txProvider repository.TxProvider,
) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		orderRepo:   orderRepo,
		packageRepo: packageRepo,
		txProvider:  txProvider,
	}
}

func (uc *DeleteOrderUseCase) Execute(ctx context.Context, id uint) error {
	// Verificar que la orden existe y obtener su estado
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// Solo se pueden eliminar órdenes pendientes
	if order.Status != "Pendiente" {
		return ErrCannotDeleteOrder
	}

	// Iniciar transacción para garantizar atomicidad
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

	// Desasociar paquetes de la orden (poner idorder = NULL)
	if err := uc.packageRepo.UnassignPackagesFromOrder(ctx, tx, id); err != nil {
		return fmt.Errorf("unassign packages from order: %w", err)
	}

	// Eliminar orden
	if err := uc.orderRepo.DeleteWithTx(ctx, tx, id); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	// Commit transacción
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}
