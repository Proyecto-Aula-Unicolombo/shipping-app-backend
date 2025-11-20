package orders

import (
	"context"
	"fmt"

	"shipping-app/internal/app/domain/ports/repository"
)

type ListOrdersInput struct {
	Limit  int
	Offset int

	TypeService string
	Status      string
	OrderID     uint
}

type ListOrdersByDriverOutput struct {
	ID          uint
	Observation *string
	Status      string
	TypeService string
}
type ListOrdersUseCase struct {
	orderRepo repository.OrderRepository
}

func NewListOrdersUseCase(orderRepo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{orderRepo: orderRepo}
}

func (uc *ListOrdersUseCase) Execute(ctx context.Context, input ListOrdersInput) ([]*ListOrdersByDriverOutput, int64, error) {
	total, err := uc.orderRepo.Count(ctx, input.TypeService, input.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}
	orders, err := uc.orderRepo.List(ctx, input.OrderID, input.Limit, input.Offset, input.TypeService, input.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}

	var listOrdersByDriverOutput []*ListOrdersByDriverOutput
	for _, order := range orders {
		listOrdersByDriverOutput = append(listOrdersByDriverOutput, &ListOrdersByDriverOutput{
			ID:          order.ID,
			Observation: order.Observation,
			Status:      order.Status,
			TypeService: order.TypeService,
		})
	}
	return listOrdersByDriverOutput, total, nil
}
