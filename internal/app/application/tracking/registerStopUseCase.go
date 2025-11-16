package tracking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/twpayne/go-geom"
)

type RegisterStopInput struct {
	OrderID     uint
	Latitude    float64
	Longitude   float64
	TypeStop    string // "Parada", "Incidente", "Recogida", "Entrega"
	Description *string
	Evidence    *string // URL de foto/evidencia
}

type RegisterStopOutput struct {
	ID        uint
	OrderID   uint
	TypeStop  string
	Timestamp time.Time
}

var (
	ErrInvalidStopInput = errors.New("invalid stop input")
	ErrInvalidStopType  = errors.New("invalid stop type: must be Parada, Incidente, Recogida, or Entrega")
	ErrOrderNotFound    = errors.New("order not found or not in progress")
)

type RegisterStopUseCase struct {
	stopRepo   repository.DeliveryStopRepository
	trackRepo  repository.TrackRepository
	orderRepo  repository.OrderRepository
	txProvider repository.TxProvider
}

func NewRegisterStopUseCase(
	stopRepo repository.DeliveryStopRepository,
	trackRepo repository.TrackRepository,
	orderRepo repository.OrderRepository,
	txProvider repository.TxProvider,
) *RegisterStopUseCase {
	return &RegisterStopUseCase{
		stopRepo:   stopRepo,
		trackRepo:  trackRepo,
		orderRepo:  orderRepo,
		txProvider: txProvider,
	}
}

func (uc *RegisterStopUseCase) Execute(ctx context.Context, input RegisterStopInput) (*RegisterStopOutput, error) {
	// Validar input
	if input.OrderID == 0 {
		return nil, ErrInvalidStopInput
	}

	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return nil, fmt.Errorf("%w: invalid coordinates", ErrInvalidStopInput)
	}

	// Validar tipo de parada
	validTypes := map[string]bool{
		"Parada":    true,
		"Incidente": true,
		"Recogida":  true,
		"Entrega":   true,
	}
	if !validTypes[input.TypeStop] {
		return nil, ErrInvalidStopType
	}

	// Verificar que la orden existe y está en progreso
	order, err := uc.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != "En camino" && order.Status != "Pendiente" {
		return nil, fmt.Errorf("order status must be 'En camino' or 'Pendiente', got: %s", order.Status)
	}

	// Iniciar transacción
	tx, err := uc.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = uc.txProvider.RollbackTx(ctx, tx)
		}
	}()

	// Crear punto geográfico
	point := geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{input.Longitude, input.Latitude})
	point.SetSRID(4326)

	now := time.Now()

	// Registrar parada de entrega
	stop := &entities.DeliveryStop{
		StopLocation: point,
		TypeStop:     input.TypeStop,
		Timestamp:    now,
		Description:  input.Description,
		Evidence:     input.Evidence,
		OrderID:      input.OrderID,
	}

	if err := uc.stopRepo.Create(ctx, tx, stop); err != nil {
		return nil, fmt.Errorf("create delivery stop: %w", err)
	}

	// También registrar en el tracking de la orden
	track := &entities.Track{
		Timestamp: now,
		Location:  point,
		OrderID:   input.OrderID,
	}

	if err := uc.trackRepo.Create(ctx, tx, track); err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}

	// Commit transacción
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	output := &RegisterStopOutput{
		ID:        stop.ID,
		OrderID:   input.OrderID,
		TypeStop:  input.TypeStop,
		Timestamp: now,
	}

	return output, nil
}
