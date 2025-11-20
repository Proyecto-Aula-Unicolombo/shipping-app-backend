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

func (r *DriverRepositoryAdapter) ListDrivers(limit, offset int, NameOrLastName string) ([]*entities.Driver, error) {
	query := `
		SELECT 
		d.id,
		u.name,
		u.lastname,
		d.phonenumber,
		d.license,
		d.is_active,
		latest_order.order_id
		FROM drivers d
		JOIN users u ON d.iduser = u.id
		LEFT JOIN  LATERAL (
			SELECT
			o.id AS order_id
			FROM orders o 
			WHERE o.iddriver = d.id
			ORDER BY o.create_at DESC
			LIMIT 1
			) AS latest_order ON true
		WHERE 1=1
	`

	args := []interface{}{}
	argsPosition := 1

	if NameOrLastName != "" {
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.lastname ILIKE $%d)", argsPosition, argsPosition)
		args = append(args, fmt.Sprintf("%%%s%%", NameOrLastName))
		argsPosition++
	}
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argsPosition, argsPosition+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*entities.Driver
	var orderID sql.NullInt64
	for rows.Next() {
		driver := entities.Driver{
			User: &entities.User{},
		}
		if err := rows.Scan(&driver.ID, &driver.User.Name, &driver.User.LastName, &driver.PhoneNumber, &driver.LicenseNo, &driver.IsActive, &orderID); err != nil {
			return nil, fmt.Errorf("error scanning driver row: %w", err)
		}
		if orderID.Valid {
			numOrder := uint(orderID.Int64)
			driver.NumOrder = numOrder
		}

		drivers = append(drivers, &driver)
	}

	return drivers, nil
}

func (r *DriverRepositoryAdapter) GetDriverByUserID(userID uint) (*entities.Driver, error) {
	query := `SELECT id, iduser, phonenumber, license FROM drivers WHERE iduser = $1`
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
	query := `
		SELECT 
		d.id,
		u.name,
		u.lastname,
		u.email,
		d.phonenumber,
		d.license,
		d.is_active,
		o.id AS order_id,
		o.status AS order_status
		FROM drivers d
		JOIN users u ON d.iduser = u.id
		LEFT JOIN orders o ON d.id = o.iddriver
		WHERE d.id = $1
		ORDER BY o.create_at DESC

	`
	driver := entities.Driver{
		User: &entities.User{},
	}
	var orderID sql.NullInt64
	var orderStatus sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&driver.ID,
		&driver.User.Name,
		&driver.User.LastName,
		&driver.User.Email,
		&driver.PhoneNumber,
		&driver.LicenseNo,
		&driver.IsActive,
		&orderID,
		&orderStatus,
	)

	if orderID.Valid {
		numOrder := uint(orderID.Int64)
		driver.NumOrder = numOrder
	}
	if orderStatus.Valid {
		driver.OrderStatus = orderStatus.String
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("driver not found with id %d", id)
		}
		return nil, fmt.Errorf("get driver by id: %w", err)
	}

	return &driver, nil
}

func (r *DriverRepositoryAdapter) CountDrivers(nameLastName string) (int64, error) {
	query := `
		SELECT COUNT(*) FROM drivers d
		JOIN users u ON d.iduser = u.id
		WHERE u.name ILIKE $1 OR u.lastname ILIKE $1`

	var count int64

	err := r.db.QueryRow(query, "%"+nameLastName+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count drivers: %w", err)
	}

	return count, nil
}

func (r *DriverRepositoryAdapter) UpdateDriverStatus(driverID uint, isActive bool) error {
	query := `UPDATE drivers SET is_active = $1 WHERE id = $2`
	_, err := r.db.Exec(query, isActive, driverID)
	if err != nil {
		return fmt.Errorf("update driver status: %w", err)
	}

	return nil
}

func (r *DriverRepositoryAdapter) ListDriversUnassigned() ([]*entities.Driver, error) {
	query := `
		SELECT 
			d.id,
			u.name,
			u.lastname,
			d.license
			FROM drivers d
			JOIN users u ON d.iduser = u.id
			LEFT JOIN orders o ON d.id = o.iddriver 
			WHERE o.id IS NULL AND d.is_active = false
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error to get drivers")
	}
	defer rows.Close()

	var drivers []*entities.Driver
	for rows.Next() {
		driver := entities.Driver{
			User: &entities.User{},
		}
		if err := rows.Scan(&driver.ID, &driver.User.Name, &driver.User.LastName, &driver.LicenseNo); err != nil {
			return nil, fmt.Errorf("error scanning driver row: %w", err)
		}

		drivers = append(drivers, &driver)
	}

	return drivers, nil
}
