package tracks

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	"time"

	"github.com/twpayne/go-geom"
)

type TrackRegisterUseCase struct {
	trackRepo repository.TrackingRepository
}

func NewTrackRegisterUseCase(trackRepo repository.TrackingRepository) *TrackRegisterUseCase {
	return &TrackRegisterUseCase{trackRepo: trackRepo}
}

type TrackRegisterInput struct {
	OrderID   uint
	Latitude  float64
	Longitude float64
}

type TrackRegisterOutput struct {
	TrackID   uint
	OrderID   uint
	Timestamp time.Time
	Longitude float64
	Latitude  float64
}

var (
	ErrGeneratePoint = errors.New("failed to generate point to location")
	ErrRegisterTrack = errors.New("failed to register track")
)

func (uc *TrackRegisterUseCase) Execute(ctx context.Context, trackInput *TrackRegisterInput) (*TrackRegisterOutput, error) {
	pointGeoLocation, err := geom.NewPoint(geom.XY).SetSRID(4326).SetCoords(geom.Coord{trackInput.Longitude, trackInput.Latitude})
	if err != nil {
		return nil, ErrGeneratePoint
	}

	track := &entities.Track{
		Location: pointGeoLocation,
		OrderID:  trackInput.OrderID,
	}

	if err := uc.trackRepo.RegisterTrack(ctx, track); err != nil {
		return nil, ErrRegisterTrack
	}

	coords := track.Location.Coords()
	return &TrackRegisterOutput{
		TrackID:   track.ID,
		OrderID:   track.OrderID,
		Timestamp: time.Now(),
		Longitude: coords.X(),
		Latitude:  coords.Y(),
	}, nil
}
