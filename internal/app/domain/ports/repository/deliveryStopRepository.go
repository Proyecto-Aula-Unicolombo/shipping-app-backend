package repository

import (
	"context"
	"database/sql"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrDeliveryStopNotFound = errors.New("delivery stop not found")
)

type DeliveryStopRepository interface {
	Create(ctx context.Context, tx *sql.Tx, stop *entities.DeliveryStop) error
	GetByID(ctx context.Context, id uint) (*entities.DeliveryStop, error)
	GetByOrderID(ctx context.Context, orderID uint) ([]*entities.DeliveryStop, error)
	ListByOrderIDWithLimit(ctx context.Context, orderID uint, limit int) ([]*entities.DeliveryStop, error)
	ListIncidents(ctx context.Context, status *string, driverID *uint, orderID *uint, limit int, offset int) ([]*entities.DeliveryStop, error)
}
