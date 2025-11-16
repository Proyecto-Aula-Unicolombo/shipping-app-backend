package orders

import (
	"context"
	"errors"
	"fmt"

	"shipping-app/internal/app/domain/ports/repository"
)

type UpdateOrderStatusInput struct {
	OrderID     uint
	Status      string
	Observation *string
}

var (
	ErrInvalidStatus = errors.New("invalid order status")
)

type UpdateOrderStatusUseCase struct {
	orderRepo repository.OrderRepository
}

func NewUpdateOrderStatusUseCase(orderRepo repository.OrderRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{orderRepo: orderRepo}
}

func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, input UpdateOrderStatusInput) error {
	// Validar que el estado es válido
	validStatuses := map[string]bool{
		"Pendiente": true,
		"En camino": true,
		"Entregado": true,
		"Cancelado": true,
	}

	if !validStatuses[input.Status] {
		return ErrInvalidStatus
	}

	// Verificar que la orden existe
	_, err := uc.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// Actualizar estado
	if err := uc.orderRepo.UpdateStatus(ctx, input.OrderID, input.Status, input.Observation); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
