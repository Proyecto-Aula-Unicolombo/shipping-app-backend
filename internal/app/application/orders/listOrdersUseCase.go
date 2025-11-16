package orders

import (
	"context"
	"fmt"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListOrdersInput struct {
	Limit  int
	Offset int
}

type ListOrdersUseCase struct {
	orderRepo repository.OrderRepository
}

func NewListOrdersUseCase(orderRepo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{orderRepo: orderRepo}
}

func (uc *ListOrdersUseCase) Execute(ctx context.Context, input ListOrdersInput) ([]*entities.Order, error) {
	orders, err := uc.orderRepo.List(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}

	return orders, nil
}
