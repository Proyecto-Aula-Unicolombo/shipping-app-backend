package repository

import (
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type DriverRepository interface {
	CreateDriverTx(tx *sql.Tx, driver *entities.Driver) error
	UpdateDriverTx(tx *sql.Tx, driver *entities.Driver) error
	ListDrivers() ([]*entities.Driver, error)
	GetDriverByUserID(userID uint) (*entities.Driver, error)
}
