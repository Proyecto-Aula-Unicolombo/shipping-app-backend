package tracks

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type GetOrderTracksUseCase struct {
	trackRepo repository.TrackRepository
	orderRepo repository.OrderRepository
}

func NewGetOrderTracksUseCase(
	trackRepo repository.TrackRepository,
	orderRepo repository.OrderRepository,
) *GetOrderTracksUseCase {
	return &GetOrderTracksUseCase{
		trackRepo: trackRepo,
		orderRepo: orderRepo,
	}
}

type GetOrderTracksInput struct {
	OrderID uint
	Limit   *int // Opcional: limitar cantidad de tracks
}

type TrackPointOutput struct {
	TrackID   uint    `json:"track_id"`
	OrderID   uint    `json:"order_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type GetOrderTracksOutput struct {
	OrderID uint               `json:"order_id"`
	Status  string             `json:"status"`
	Tracks  []TrackPointOutput `json:"tracks"`
}

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrNoTracksFound = errors.New("no tracks found for this order")
)

func (uc *GetOrderTracksUseCase) Execute(ctx context.Context, input GetOrderTracksInput) (*GetOrderTracksOutput, error) {
	// Verificar que la orden existe
	order, err := uc.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Obtener tracks
	var tracks []*entities.Track
	if input.Limit != nil && *input.Limit > 0 {
		tracks, err = uc.trackRepo.ListByOrderIDWithLimit(ctx, input.OrderID, *input.Limit)
	} else {
		tracks, err = uc.trackRepo.GetByOrderID(ctx, input.OrderID)
	}

	if err != nil {
		if errors.Is(err, repository.ErrTrackNotFound) {
			return nil, ErrNoTracksFound
		}
		return nil, err
	}

	// Convertir a output
	trackPoints := make([]TrackPointOutput, 0, len(tracks))
	for _, track := range tracks {
		coords := track.Location.Coords()
		trackPoints = append(trackPoints, TrackPointOutput{
			TrackID:   track.ID,
			OrderID:   track.OrderID,
			Latitude:  coords.Y(),
			Longitude: coords.X(),
			Timestamp: track.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &GetOrderTracksOutput{
		OrderID: order.ID,
		Status:  order.Status,
		Tracks:  trackPoints,
	}, nil
}
