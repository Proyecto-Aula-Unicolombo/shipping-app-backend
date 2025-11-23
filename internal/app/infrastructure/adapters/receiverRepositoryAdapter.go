package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
)

type ReceiverRepositoryPostgres struct {
	db *sql.DB
}

func NewReceiverRepositoryPostgres(db *sql.DB) *ReceiverRepositoryPostgres {
	return &ReceiverRepositoryPostgres{db: db}
}

func (r *ReceiverRepositoryPostgres) GetByID(ctx context.Context, id uint) (*entities.Receiver, error) {
	query := `SELECT id, name, lastname, phonenumber, email FROM receivers WHERE id = $1`
	var rc entities.Receiver
	var row *sql.Row
	row = r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&rc.ID, &rc.Name, &rc.LastName, &rc.PhoneNumber, &rc.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("receiver get: %w", err)
	}
	return &rc, nil
}

func (r *ReceiverRepositoryPostgres) FindByEmail(ctx context.Context, email string) (*entities.Receiver, error) {
	return r.FindByEmailWithTx(ctx, nil, email)
}

func (r *ReceiverRepositoryPostgres) FindByEmailWithTx(ctx context.Context, tx *sql.Tx, email string) (*entities.Receiver, error) {
	query := `
		SELECT id, name, lastname, phonenumber, email
		FROM receivers
		WHERE email = $1
		LIMIT 1
	`

	var receiver entities.Receiver
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, email).Scan(
			&receiver.ID,
			&receiver.Name,
			&receiver.LastName,
			&receiver.PhoneNumber,
			&receiver.Email,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, email).Scan(
			&receiver.ID,
			&receiver.Name,
			&receiver.LastName,
			&receiver.PhoneNumber,
			&receiver.Email,
		)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find receiver by email: %w", err)
	}

	return &receiver, nil
}

func (r *ReceiverRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, receiver *entities.Receiver) error {
	query := `
		INSERT INTO receivers (name, lastname, phonenumber, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	args := []interface{}{
		receiver.Name,
		receiver.LastName,
		receiver.PhoneNumber,
		receiver.Email,
	}

	var err error
	if tx != nil {
		log.Println("Creating Receiver within transaction")
		err = tx.QueryRowContext(ctx, query, args...).Scan(&receiver.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&receiver.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("receiver create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("receiver create timeout: %w", err)
		}
		return fmt.Errorf("receiver create: %w", err)
	}

	log.Printf("Receiver created successfully with ID: %d", receiver.ID)
	return nil
}
