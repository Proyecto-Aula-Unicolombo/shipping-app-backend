package orders

import (
	"context"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListOrdersUnassignedInput struct {
	Limit, Offset int
	ID            uint
}

type ListOrdersUnassignedOutput struct {
	ID          uint
	Status      string
	Typeservice string
}

type ListOrdersUnassignedUseCase struct {
	orderRepository repository.OrderRepository
}

func NewListOrdersUnassignedUseCase(orderRepository repository.OrderRepository) *ListOrdersUnassignedUseCase {
	return &ListOrdersUnassignedUseCase{orderRepository: orderRepository}
}

func (uc *ListOrdersUnassignedUseCase) Execute(ctx context.Context, input ListOrdersUnassignedInput) ([]ListOrdersUnassignedOutput, int64, error) {
	orders, err := uc.orderRepository.ListOrderUnassigned(ctx, input.Limit, input.Offset, input.ID)
	if err != nil {
		return nil, 0, err
	}
	total, err := uc.orderRepository.Count(ctx, "", "")
	if err != nil {
		return nil, 0, err
	}

	var outputs []ListOrdersUnassignedOutput
	for _, order := range orders {
		outputs = append(outputs, ListOrdersUnassignedOutput{
			ID:          order.ID,
			Status:      order.Status,
			Typeservice: order.TypeService,
		})
	}

	return outputs, total, nil
}
