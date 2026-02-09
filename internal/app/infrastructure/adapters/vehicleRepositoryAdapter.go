package adapters

import (
	"context"
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

func (r *VehicleRepositoryPostgres) CreateVehicleTx(v *entities.Vehicle) error {
	query := `
		INSERT INTO vehicles (plate, brand, model, color, vehicletype)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var err error

	err = r.db.QueryRow(query,
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

func (r *VehicleRepositoryPostgres) GetByID(ctx context.Context, id uint) (*entities.Vehicle, error) {
	var v entities.Vehicle

	query := `
		SELECT id, plate, brand, model, color, vehicletype
		FROM vehicles
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&v.ID,
		&v.Plate,
		&v.Brand,
		&v.Model,
		&v.Color,
		&v.VehicleType,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVehicleNotFound
		}
		return nil, fmt.Errorf("get vehicle by id: %w", err)
	}

	return &v, nil
}

func (r *VehicleRepositoryPostgres) HasActiveVehicleInOrder(ctx context.Context, id uint) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM orders 
			WHERE idvehicle = $1
			AND status IN ('asignada', 'en camino')
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
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

func (r *VehicleRepositoryPostgres) ListVehicles(limit int, offset int, PlateBrandOrModel string) ([]*entities.Vehicle, error) {
	query := `
		SELECT 
            v.id, 
            v.plate, 
            v.brand, 
            v.model,
			v.color,
            v.vehicletype,
            u.name AS driver_name,
            u.lastname AS driver_last_name
        FROM vehicles v
        LEFT JOIN orders o ON v.id = o.idvehicle 
        LEFT JOIN drivers d ON o.iddriver = d.id
        LEFT JOIN users u ON d.iduser = u.id
        WHERE 1=1
	`
	args := []interface{}{}
	argPosition := 1

	if PlateBrandOrModel != "" {
		query += fmt.Sprintf(" AND (plate ILIKE $%d OR brand ILIKE $%d OR model ILIKE $%d)", argPosition, argPosition, argPosition)
		args = append(args, "%"+PlateBrandOrModel+"%")
		argPosition++
	}

	query += " ORDER BY id LIMIT $" + fmt.Sprint(argPosition) + " OFFSET $" + fmt.Sprint(argPosition+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer rows.Close()

	var vehicles []*entities.Vehicle
	var driverName sql.NullString
	var driverLastName sql.NullString
	for rows.Next() {
		var vehicle entities.Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Plate, &vehicle.Brand, &vehicle.Model, &vehicle.Color, &vehicle.VehicleType, &driverName, &driverLastName); err != nil {
			return nil, fmt.Errorf("error scanning vehicle row: %w", err)
		}
		if driverName.Valid {
			vehicle.AssignedDriverName = driverName.String
		}
		if driverLastName.Valid {
			vehicle.AssignedDriverLastName = driverLastName.String
		}
		vehicles = append(vehicles, &vehicle)
	}
	return vehicles, nil
}

func (r *VehicleRepositoryPostgres) CountVehicles(PlateBrandOrModel string) (int64, error) {
	query := `SELECT COUNT(*) FROM vehicles WHERE (plate ILIKE $1 OR brand ILIKE $1 OR model ILIKE $1)`
	var count int64
	err := r.db.QueryRow(query, "%"+PlateBrandOrModel+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting vehicles: %w", err)
	}
	return count, nil
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

func (r *VehicleRepositoryPostgres) ListVehiclesUnassigned() ([]*entities.Vehicle, error) {
	query := `
		SELECT
			v.id,
			v.plate,
			v.brand,
			v.model,
			v.vehicletype
		FROM vehicles v
		LEFT JOIN orders o ON v.id = o.idvehicle
		WHERE o.id IS NULL
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listing unassigned vehicles: %w", err)
	}

	defer rows.Close()

	var vehicles []*entities.Vehicle
	for rows.Next() {
		var vehicle entities.Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Plate, &vehicle.Brand, &vehicle.Model, &vehicle.VehicleType); err != nil {
			return nil, fmt.Errorf("error scanning vehicle row: %w", err)
		}
		vehicles = append(vehicles, &vehicle)
	}
	return vehicles, nil
}
