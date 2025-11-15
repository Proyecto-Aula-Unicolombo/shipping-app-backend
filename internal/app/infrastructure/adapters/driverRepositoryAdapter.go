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

func (r *DriverRepositoryAdapter) ListDrivers() ([]*entities.Driver, error) {
	query := `SELECT id, iduser, phonenumber, license FROM drivers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var drivers []*entities.Driver
	for rows.Next() {
		var driver entities.Driver
		if err := rows.Scan(&driver.ID, &driver.UserID, &driver.PhoneNumber, &driver.LicenseNo); err != nil {
			return nil, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, nil
}

func (r *DriverRepositoryAdapter) GetDriverByUserID(userID uint) (*entities.Driver, error) {
	query := `SELECT id, phonenumber, license FROM drivers WHERE iduser = $1`
	row := r.db.QueryRow(query, userID)
	var driver entities.Driver
	if err := row.Scan(&driver.ID, &driver.PhoneNumber, &driver.LicenseNo); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepositoryAdapter) UpdateDriverTx(tx *sql.Tx, driver *entities.Driver) error {
	query := `UPDATE drivers SET phonenumber = $1, license = $2 WHERE iduser = $3`

	_, err := tx.Exec(query, driver.PhoneNumber, driver.LicenseNo, driver.UserID)
	return err

}

func (r *DriverRepositoryAdapter) DeleteDriverByUserIDTx(tx *sql.Tx, userID uint) error {
	query := `DELETE FROM drivers WHERE iduser = $1`
	_, err := tx.Exec(query, userID)
	return err
}
