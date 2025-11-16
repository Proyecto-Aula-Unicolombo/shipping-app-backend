package repository

import (
	"context"
	"database/sql"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrTrackNotFound = errors.New("track not found")
)

type TrackRepository interface {
	Create(ctx context.Context, tx *sql.Tx, track *entities.Track) error
	GetByOrderID(ctx context.Context, orderID uint) ([]*entities.Track, error)
	GetLatestByOrderID(ctx context.Context, orderID uint) (*entities.Track, error)
	ListByOrderIDWithLimit(ctx context.Context, orderID uint, limit int) ([]*entities.Track, error)
}
