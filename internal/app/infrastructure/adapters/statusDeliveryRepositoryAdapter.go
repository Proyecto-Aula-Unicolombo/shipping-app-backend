package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
)

type StatusDeliveryRepositoryPostgres struct {
	db *sql.DB
}

func NewStatusDeliveryRepositoryPostgres(db *sql.DB) *StatusDeliveryRepositoryPostgres {
	return &StatusDeliveryRepositoryPostgres{db: db}
}

func (r *StatusDeliveryRepositoryPostgres) GetByID(ctx context.Context, tx *sql.Tx, id uint) (*entities.StatusDelivery, error) {
	query := `SELECT id, status, priority, date_estimated_delivery, date_real_delivery FROM statusdelivery WHERE id = $1`
	var s entities.StatusDelivery
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	} else {
		row = r.db.QueryRowContext(ctx, query, id)
	}
	if err := row.Scan(&s.ID, &s.Status, &s.Priority, &s.DateEstimatedDelivery, &s.DateRealDelivery); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("status get: %w", err)
	}
	return &s, nil
}

func (r *StatusDeliveryRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, status *entities.StatusDelivery) error {
	query := `
		INSERT INTO statusdelivery (status, priority, date_estimated_delivery, date_real_delivery)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	args := []interface{}{
		status.Status,
		status.Priority,
		status.DateEstimatedDelivery,
		status.DateRealDelivery,
	}

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&status.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&status.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("status delivery create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("status delivery create timeout: %w", err)
		}
		return fmt.Errorf("status delivery create: %w", err)
	}

	return nil
}

func (r *StatusDeliveryRepositoryPostgres) Delete(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `DELETE FROM statusdelivery WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, id)
	}

	if err != nil {
		log.Printf("ERROR executing DELETE statusdelivery: %v", err)
		return fmt.Errorf("delete status delivery: %w", err)
	}

	return nil

}
