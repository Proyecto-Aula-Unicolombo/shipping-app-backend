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

// ==================== CREATE ====================

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

// ==================== GET BY ID ====================

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

// ==================== UNIMPLEMENTED ====================

func (r *VehicleRepositoryPostgres) DeleteVehicle(id uint) error {
	panic("unimplemented")
}

func (r *VehicleRepositoryPostgres) GetAllVehicles() ([]*entities.Vehicle, error) {
	panic("unimplemented")
}

func (r *VehicleRepositoryPostgres) ListVehicles(limit int, offset int, PlateOrBrand string) ([]*entities.Vehicle, error) {
	panic("unimplemented")
}

func (r *VehicleRepositoryPostgres) UpdateVehicle(vehicle *entities.Vehicle) error {
	panic("unimplemented")
}
