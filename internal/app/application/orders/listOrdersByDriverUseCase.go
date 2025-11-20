package orders

import (
	"context"
	"fmt"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListOrdersByDriverInput struct {
	DriverID uint
	Limit    int
	Offset   int
}

type ListOrdersByDriverUseCase struct {
	orderRepo repository.OrderRepository
}

func NewListOrdersByDriverUseCase(orderRepo repository.OrderRepository) *ListOrdersByDriverUseCase {
	return &ListOrdersByDriverUseCase{orderRepo: orderRepo}
}

func (uc *ListOrdersByDriverUseCase) Execute(ctx context.Context, input ListOrdersByDriverInput) ([]*entities.Order, int64, error) {
	total, err := uc.orderRepo.CountByDriver(ctx, input.DriverID)
	if err != nil {
		return nil, 0, fmt.Errorf("count orders by driver: %w", err)
	}

	orders, err := uc.orderRepo.ListByDriver(ctx, input.DriverID, input.Limit, input.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders by driver: %w", err)
	}

	return orders, total, nil
}
