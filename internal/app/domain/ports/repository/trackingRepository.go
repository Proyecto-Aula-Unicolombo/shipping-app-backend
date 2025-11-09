package repository

import (
	"context"
	"shipping-app/internal/app/domain/entities"
)

type TrackingRepository interface {
	RegisterTrack(ctx context.Context, track *entities.Track) error
}
