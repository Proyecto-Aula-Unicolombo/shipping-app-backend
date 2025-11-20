package repository

import (
	"context"
	"database/sql"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository interface {
	Create(ctx context.Context, tx *sql.Tx, order *entities.Order) error
	Update(ctx context.Context, order *entities.Order) error
	UpdateStatus(ctx context.Context, id uint, status string, observation *string) error
	AssignDriverAndVehicle(ctx context.Context, id uint, driverID, vehicleID uint) error
	Delete(ctx context.Context, id uint) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id uint) error
	GetByID(ctx context.Context, id uint) (*entities.Order, error)
	List(ctx context.Context, orderID uint, limit, offset int, typeService, status string) ([]*entities.Order, error)
	ListByDriver(ctx context.Context, driverID uint, limit, offset int) ([]*entities.Order, error)
	Count(ctx context.Context, typeService, status string) (int64, error)
	CountByDriver(ctx context.Context, driverID uint) (int64, error)
}
