package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shipping-app/internal/app/domain/entities"
)

type AddressPackageRepositoryPostgres struct {
	db *sql.DB
}

func NewAddressPackageRepositoryPostgres(db *sql.DB) *AddressPackageRepositoryPostgres {
	return &AddressPackageRepositoryPostgres{db: db}
}

func (r *AddressPackageRepositoryPostgres) GetByID(ctx context.Context, tx *sql.Tx, id uint) (*entities.AddressPackage, error) {
	query := `SELECT id, origin, destination, delivery_instructions FROM addresspackages WHERE id = $1`
	var addr entities.AddressPackage
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	} else {
		row = r.db.QueryRowContext(ctx, query, id)
	}
	if err := row.Scan(&addr.ID, &addr.Origin, &addr.Destination, &addr.DeliveryInstructions); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no existe
		}
		return nil, fmt.Errorf("address get: %w", err)
	}
	return &addr, nil
}
func (r *AddressPackageRepositoryPostgres) FindByRoute(ctx context.Context, tx *sql.Tx, origin, destination string) (*entities.AddressPackage, error) {
	query := `
		SELECT id, origin, destination, delivery_instructions
		FROM addresspackages
		WHERE origin = $1 AND destination = $2
		LIMIT 1
	`

	var addr entities.AddressPackage
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, origin, destination).Scan(
			&addr.ID,
			&addr.Origin,
			&addr.Destination,
			&addr.DeliveryInstructions,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, origin, destination).Scan(
			&addr.ID,
			&addr.Origin,
			&addr.Destination,
			&addr.DeliveryInstructions,
		)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No encontrado, no es error
		}
		return nil, fmt.Errorf("find address by route: %w", err)
	}

	return &addr, nil
}

func (r *AddressPackageRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, addr *entities.AddressPackage) error {
	query := `
		INSERT INTO addresspackages (origin, destination, delivery_instructions)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []interface{}{addr.Origin, addr.Destination, addr.DeliveryInstructions}

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&addr.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&addr.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("address create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("address create timeout: %w", err)
		}
		return fmt.Errorf("address create: %w", err)
	}

	return nil
}
