package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
)

type SenderRepositoryPostgres struct {
	db *sql.DB
}

func NewSenderRepositoryPostgres(db *sql.DB) *SenderRepositoryPostgres {
	return &SenderRepositoryPostgres{db: db}
}

func (r *SenderRepositoryPostgres) GetByID(ctx context.Context, tx *sql.Tx, id uint) (*entities.Sender, error) {
	query := `SELECT id, name, document, address, phonenumber, email FROM senders WHERE id = $1`
	var s entities.Sender
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	} else {
		row = r.db.QueryRowContext(ctx, query, id)
	}
	if err := row.Scan(&s.ID, &s.Name, &s.Document, &s.Address, &s.PhoneNumber, &s.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("sender get: %w", err)
	}
	return &s, nil
}

func (r *SenderRepositoryPostgres) FindByEmailOrDocument(ctx context.Context, tx *sql.Tx, email, document string) (*entities.Sender, error) {
	query := `
		SELECT id, name, document, address, phonenumber, email
		FROM senders
		WHERE email = $1 OR document = $2
		LIMIT 1
	`

	var sender entities.Sender
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, email, document).Scan(
			&sender.ID,
			&sender.Name,
			&sender.Document,
			&sender.Address,
			&sender.PhoneNumber,
			&sender.Email,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, email, document).Scan(
			&sender.ID,
			&sender.Name,
			&sender.Document,
			&sender.Address,
			&sender.PhoneNumber,
			&sender.Email,
		)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find sender by email or document: %w", err)
	}

	return &sender, nil
}

func (r *SenderRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, sender *entities.Sender) error {
	query := `
		INSERT INTO senders (name, document, address, phonenumber, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	args := []interface{}{
		sender.Name,
		sender.Document,
		sender.Address,
		sender.PhoneNumber,
		sender.Email,
	}

	var err error
	if tx != nil {
		log.Println("Creating Sender within transaction")
		err = tx.QueryRowContext(ctx, query, args...).Scan(&sender.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&sender.ID)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("sender create canceled: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("sender create timeout: %w", err)
		}
		return fmt.Errorf("sender create: %w", err)
	}

	log.Printf("Sender created successfully with ID: %d", sender.ID)
	return nil
}
