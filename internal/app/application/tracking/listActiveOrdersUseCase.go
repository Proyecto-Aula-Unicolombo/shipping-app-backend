package tracking

import (
	"context"
	"shipping-app/internal/app/domain/ports/repository"
)

type OrderLocationOutput struct {
	OrderID    uint     `json:"order_id"`
	Status     string   `json:"status"`
	DriverID   *uint    `json:"driver_id,omitempty"`
	DriverName *string  `json:"driver_name,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Timestamp  *string  `json:"timestamp,omitempty"`
}

type ListActiveOrdersUseCase struct {
	orderRepo repository.OrderRepository
	trackRepo repository.TrackRepository
}

func NewListActiveOrdersUseCase(
	orderRepo repository.OrderRepository,
	trackRepo repository.TrackRepository,
) *ListActiveOrdersUseCase {
	return &ListActiveOrdersUseCase{
		orderRepo: orderRepo,
		trackRepo: trackRepo,
	}
}

func (uc *ListActiveOrdersUseCase) Execute(ctx context.Context) ([]OrderLocationOutput, error) {
	// Obtener todas las órdenes activas (En camino o Pendiente)
	// Usar List con límite alto para obtener todas
	orders, err := uc.orderRepo.List(ctx, 0, 1000, 0, "", "")
	if err != nil {
		return nil, err
	}

	results := make([]OrderLocationOutput, 0)

	for _, order := range orders {
		// Solo incluir órdenes "En camino" o "Pendiente"
		if order.Status != "en camino" && order.Status != "pendiente" {
			continue
		}

		output := OrderLocationOutput{
			OrderID:  order.ID,
			Status:   order.Status,
			DriverID: order.DriverID,
		}

		// Obtener nombre del conductor si está asignado
		if order.Driver != nil && order.Driver.User != nil {
			if order.Driver.User.Name != "" {
				driverName := order.Driver.User.Name
				if order.Driver.User.LastName != "" {
					driverName += " " + order.Driver.User.LastName
				}
				output.DriverName = &driverName
			}
		}

		// Obtener el último track de esta orden
		latestTrack, err := uc.trackRepo.GetLatestByOrderID(ctx, order.ID)
		if err == nil && latestTrack != nil {
			lat := latestTrack.Location.Y()
			lng := latestTrack.Location.X()
			timestamp := latestTrack.Timestamp.Format("2006-01-02T15:04:05Z07:00")

			output.Latitude = &lat
			output.Longitude = &lng
			output.Timestamp = &timestamp
		}

		results = append(results, output)
	}

	return results, nil
}
