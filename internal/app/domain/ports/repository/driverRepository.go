package repository

import (
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type DriverRepository interface {
	CreateDriverTx(tx *sql.Tx, driver *entities.Driver) error
}
