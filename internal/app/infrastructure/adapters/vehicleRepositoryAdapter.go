package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
	"github.com/lib/pq"
)

var (
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
)

type VehicleRepositoryPostgres struct {
	db *sql.DB
}

func NewVehicleRepositoryPostgres(db *sql.DB) *VehicleRepositoryPostgres {
	return &VehicleRepositoryPostgres{db: db}
}


func (r *VehicleRepositoryPostgres) CreateVehicleTx(tx *sql.Tx, v *entities.Vehicle) error {
	query := `
		INSERT INTO vehicles (plate, brand, model, color, vehicletype)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var err error

	err = tx.QueryRow(query,
		v.Plate,
		v.Brand,
		v.Model,
		v.Color,
		v.VehicleType,
	).Scan(&v.ID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrVehicleAlreadyExists
			}
		}
		return fmt.Errorf("error creating vehicle: %w", err)
	}

	return nil
}


func (r *VehicleRepositoryPostgres) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	var v entities.Vehicle

	query := `
		SELECT id, plate, brand, model, color, vehicletype
		FROM vehicles
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&v.ID,
		&v.Plate,
		&v.Brand,
		&v.Model,
		&v.Color,
		&v.VehicleType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrVehicleNotFound
		}
		return nil, fmt.Errorf("error fetching vehicle by id: %w", err)
	}

	return &v, nil
}


func (r *VehicleRepositoryPostgres) DeleteVehicle(id uint) error {
	query := `DELETE FROM vehicles WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting vehicle: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVehicleNotFound
	}

	return nil
}

func (r *VehicleRepositoryPostgres) GetAllVehicles() ([]*entities.Vehicle, error) {
	query := `
		SELECT id, plate, brand, model, color, vehicletype
		FROM vehicles
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting all vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*entities.Vehicle
	for rows.Next() {
		var v entities.Vehicle
		err := rows.Scan(
			&v.ID,
			&v.Plate,
			&v.Brand,
			&v.Model,
			&v.Color,
			&v.VehicleType,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}
		vehicles = append(vehicles, &v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating vehicles: %w", err)
	}

	return vehicles, nil
}

func (r *VehicleRepositoryPostgres) ListVehicles(limit int, offset int, PlateOrBrand string) ([]*entities.Vehicle, error) {
	panic("unimplemented")
}

func (r *VehicleRepositoryPostgres) UpdateVehicle(vehicle *entities.Vehicle) error {
	query := `
		UPDATE vehicles 
		SET plate = $1, brand = $2, model = $3, color = $4, vehicletype = $5
		WHERE id = $6
	`

	res, err := r.db.Exec(
		query,
		vehicle.Plate,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Color,
		vehicle.VehicleType,
		vehicle.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating vehicle: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVehicleNotFound
	}

	return nil
}
