package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
)

type ComercialInformationRepositoryPostgres struct {
	db *sql.DB
}

func NewComercialInformationRepositoryPostgres(db *sql.DB) *ComercialInformationRepositoryPostgres {
	return &ComercialInformationRepositoryPostgres{db: db}
}

func (r *ComercialInformationRepositoryPostgres) GetByID(ctx context.Context, tx *sql.Tx, id uint) (*entities.ComercialInformation, error) {
	query := `SELECT id, cost_sending, is_paid FROM comercialinformations WHERE id = $1`
	var ci entities.ComercialInformation
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	} else {
		row = r.db.QueryRowContext(ctx, query, id)
	}
	if err := row.Scan(&ci.ID, &ci.CostSending, &ci.IsPaid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("comercialinfo get: %w", err)
	}
	return &ci, nil
}

func (r *ComercialInformationRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, ci *entities.ComercialInformation) error {
	query := `
		INSERT INTO comercialinformations (cost_sending, is_paid)
		VALUES ($1, $2)
		RETURNING id
	`

	args := []interface{}{ci.CostSending, ci.IsPaid}

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&ci.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&ci.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("comercialinfo create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("comercialinfo create timeout: %w", err)
		}
		return fmt.Errorf("comercialinfo create: %w", err)
	}

	return nil
}

func (r *ComercialInformationRepositoryPostgres) Delete(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `DELETE FROM comercialinformations WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, id)
	}

	if err != nil {
		log.Printf("ERROR executing DELETE comercialinformations: %v", err)
		return fmt.Errorf("delete comercial information: %w", err)
	}

	return nil
}
