package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (r *DriverRepositoryAdapter) UpdateDriverTx(tx *sql.Tx, driver *entities.Driver) error {
	query := `UPDATE drivers SET phonenumber = $1, license = $2 WHERE iduser = $3`
	if tx != nil {
		_, err := tx.Exec(query, driver.PhoneNumber, driver.LicenseNo, driver.UserID)
		return err
	}
	_, err := r.db.Exec(query, driver.PhoneNumber, driver.LicenseNo, driver.UserID)
	return err
}

func (r *DriverRepositoryAdapter) ListDrivers() ([]*entities.Driver, error) {
	query := `SELECT iddriver, iduser, phonenumber, license FROM drivers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []*entities.Driver
	for rows.Next() {
		driver := &entities.Driver{}
		err := rows.Scan(&driver.ID, &driver.UserID, &driver.PhoneNumber, &driver.LicenseNo)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (r *DriverRepositoryAdapter) GetDriverByUserID(userID uint) (*entities.Driver, error) {
	query := `SELECT iddriver, iduser, phonenumber, license FROM drivers WHERE iduser = $1`
	driver := &entities.Driver{}
	err := r.db.QueryRow(query, userID).Scan(&driver.ID, &driver.UserID, &driver.PhoneNumber, &driver.LicenseNo)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (r *DriverRepositoryAdapter) DeleteDriverByUserIDTx(tx *sql.Tx, userID uint) error {
	query := `DELETE FROM drivers WHERE iduser = $1`
	if tx != nil {
		_, err := tx.Exec(query, userID)
		return err
	}
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *DriverRepositoryAdapter) GetByID(ctx context.Context, id uint) (*entities.Driver, error) {
	query := `SELECT id, phonenumber, license, iduser FROM drivers WHERE id = $1`

	var driver entities.Driver
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&driver.ID,
		&driver.PhoneNumber,
		&driver.LicenseNo,
		&driver.UserID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("driver not found with id %d", id)
		}
		return nil, fmt.Errorf("get driver by id: %w", err)
	}

	return &driver, nil
}
