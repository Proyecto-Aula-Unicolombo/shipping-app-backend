package tracking

import (
	"context"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListIncidentsInput struct {
	Status   *string // "Abierto", "En Progreso", "Cerrado"
	DriverID *uint
	OrderID  *uint
	Limit    int
	Offset   int
}

type IncidentOutput struct {
	ID          uint    `json:"id"`
	OrderID     uint    `json:"order_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TypeStop    string  `json:"type_stop"`
	Description *string `json:"description"`
	Evidence    *string `json:"evidence"`
	Timestamp   string  `json:"timestamp"`
	Status      string  `json:"status"`
	// Driver info
	DriverID       *uint   `json:"driver_id,omitempty"`
	DriverName     *string `json:"driver_name,omitempty"`
	DriverLastName *string `json:"driver_lastname,omitempty"`
}

type ListIncidentsUseCase struct {
	stopRepo  repository.DeliveryStopRepository
	orderRepo repository.OrderRepository
}

func NewListIncidentsUseCase(
	stopRepo repository.DeliveryStopRepository,
	orderRepo repository.OrderRepository,
) *ListIncidentsUseCase {
	return &ListIncidentsUseCase{
		stopRepo:  stopRepo,
		orderRepo: orderRepo,
	}
}

func (uc *ListIncidentsUseCase) Execute(ctx context.Context, input ListIncidentsInput) ([]IncidentOutput, error) {
	// Por ahora, obtenemos todos los delivery stops tipo "Incidente"
	// TODO: Implementar filtros en el repositorio

	stops, err := uc.stopRepo.ListIncidents(ctx, input.Status, input.DriverID, input.OrderID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	incidents := make([]IncidentOutput, 0, len(stops))
	for _, stop := range stops {
		incident := IncidentOutput{
			ID:          stop.ID,
			OrderID:     stop.OrderID,
			Latitude:    stop.StopLocation.Y(),
			Longitude:   stop.StopLocation.X(),
			TypeStop:    stop.TypeStop,
			Description: stop.Description,
			Evidence:    stop.Evidence,
			Timestamp:   stop.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Determinar el estado basado en la evidencia o timestamp
		if stop.Evidence != nil && *stop.Evidence != "" {
			incident.Status = "Cerrado"
		} else {
			// Si es reciente (menos de 1 hora), está "Abierto", sino "En Progreso"
			// Este es un ejemplo, puedes ajustar la lógica
			incident.Status = "Abierto"
		}

		// Obtener información del conductor si está disponible
		if stop.Order != nil && stop.Order.DriverID != nil {
			incident.DriverID = stop.Order.DriverID
			if stop.Order.Driver != nil {
				if stop.Order.Driver.User.Name != "" {
					incident.DriverName = &stop.Order.Driver.User.Name
				}
				if stop.Order.Driver.User.LastName != "" {
					incident.DriverLastName = &stop.Order.Driver.User.LastName
				}
			}
		}

		incidents = append(incidents, incident)
	}

	return incidents, nil
}
