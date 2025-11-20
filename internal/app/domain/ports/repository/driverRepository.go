package repository

import (
	"context"
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type DriverRepository interface {
	CreateDriverTx(tx *sql.Tx, driver *entities.Driver) error
	UpdateDriverTx(tx *sql.Tx, driver *entities.Driver) error
	ListDrivers(limit, offset int, nameLastNameOrNumOrder string) ([]*entities.Driver, error)
	GetDriverByUserID(userID uint) (*entities.Driver, error)
	DeleteDriverByUserIDTx(tx *sql.Tx, userID uint) error
	GetByID(ctx context.Context, id uint) (*entities.Driver, error)
	CountDrivers(nameLastNameOrNumOrder string) (int64, error)
	UpdateDriverStatus(driverID uint, isActive bool) error
	ListDriversUnassigned() ([]*entities.Driver, error)
}
