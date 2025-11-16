package orders

import (
	"context"
	"fmt"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type GetOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewGetOrderUseCase(orderRepo repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, id uint) (*entities.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return order, nil
}
