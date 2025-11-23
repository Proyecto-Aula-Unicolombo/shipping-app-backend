package tracks

import (
	"context"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type GetActiveDriversLocationsUseCase struct {
	trackRepo  repository.TrackRepository
	orderRepo  repository.OrderRepository
	driverRepo repository.DriverRepository
}

func NewGetActiveDriversLocationsUseCase(
	trackRepo repository.TrackRepository,
	orderRepo repository.OrderRepository,
	driverRepo repository.DriverRepository,
) *GetActiveDriversLocationsUseCase {
	return &GetActiveDriversLocationsUseCase{
		trackRepo:  trackRepo,
		orderRepo:  orderRepo,
		driverRepo: driverRepo,
	}
}

type DriverLocationOutput struct {
	DriverID     uint    `json:"driver_id"`
	DriverName   string  `json:"driver_name"`
	PhoneNumber  string  `json:"phone_number"`
	OrderID      uint    `json:"order_id"`
	OrderStatus  string  `json:"order_status"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	LastUpdate   string  `json:"last_update"`
	VehicleID    *uint   `json:"vehicle_id,omitempty"`
	VehiclePlate string  `json:"vehicle_plate,omitempty"`
}

type GetActiveDriversLocationsOutput struct {
	TotalDrivers int                    `json:"total_drivers"`
	Drivers      []DriverLocationOutput `json:"drivers"`
}

func (uc *GetActiveDriversLocationsUseCase) Execute(ctx context.Context) (*GetActiveDriversLocationsOutput, error) {
	// Obtener todas las órdenes activas (no entregadas ni canceladas)
	// Por ahora usamos List sin filtros y filtramos manualmente
	// TODO: Agregar método específico al repositorio para órdenes activas
	orders, err := uc.orderRepo.List(ctx, 0, 1000, 0, "", "")
	if err != nil {
		return nil, err
	}

	// Filtrar órdenes activas con conductor asignado
	activeOrders := make([]*entities.Order, 0)
	for _, order := range orders {
		if order.DriverID != nil &&
			order.Status != "Entregado" &&
			order.Status != "Cancelado" {
			activeOrders = append(activeOrders, order)
		}
	}

	// Para cada orden activa, obtener la última ubicación
	driverLocations := make([]DriverLocationOutput, 0)
	driverMap := make(map[uint]bool) // Para evitar duplicados

	for _, order := range activeOrders {
		// Si ya procesamos este conductor, continuar
		if driverMap[*order.DriverID] {
			continue
		}

		// Obtener última ubicación del conductor para esta orden
		track, err := uc.trackRepo.GetLatestByOrderID(ctx, order.ID)
		if err != nil {
			// Si no hay tracks, continuar con la siguiente orden
			continue
		}

		// Obtener información del conductor
		driver, err := uc.driverRepo.GetByID(ctx, *order.DriverID)
		if err != nil {
			continue
		}

		coords := track.Location.Coords()
		driverName := ""
		if driver.User != nil {
			driverName = driver.User.Name + " " + driver.User.LastName
		}

		vehiclePlate := ""
		if order.Vehicle != nil {
			vehiclePlate = order.Vehicle.Plate
		}

		driverLocation := DriverLocationOutput{
			DriverID:     *order.DriverID,
			DriverName:   driverName,
			PhoneNumber:  driver.PhoneNumber,
			OrderID:      order.ID,
			OrderStatus:  order.Status,
			Latitude:     coords.Y(),
			Longitude:    coords.X(),
			LastUpdate:   track.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			VehicleID:    order.VehicleID,
			VehiclePlate: vehiclePlate,
		}

		driverLocations = append(driverLocations, driverLocation)
		driverMap[*order.DriverID] = true
	}

	return &GetActiveDriversLocationsOutput{
		TotalDrivers: len(driverLocations),
		Drivers:      driverLocations,
	}, nil
}
