package adapters

import (
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type DriverRepositoryAdapter struct {
	db *sql.DB
}

func NewDriverRepositoryAdapter(db *sql.DB) *DriverRepositoryAdapter {
	return &DriverRepositoryAdapter{db: db}
}

func (r *DriverRepositoryAdapter) CreateDriverTx(tx *sql.Tx, driver *entities.Driver) error {
	query := `INSERT INTO drivers (iduser, phonenumber, license) VALUES ($1, $2, $3)`
	if tx != nil {
		_, err := tx.Exec(query, driver.UserID, driver.PhoneNumber, driver.LicenseNo)
		return err
	}
	_, err := r.db.Exec(query, driver.UserID, driver.PhoneNumber, driver.LicenseNo)
	return err
}
