package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type InformationDeliveryRepositoryPostgres struct {
	db *sql.DB
}

func NewInformationDeliveryRepositoryPostgres(db *sql.DB) *InformationDeliveryRepositoryPostgres {
	return &InformationDeliveryRepositoryPostgres{db: db}
}

func (r *InformationDeliveryRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, info *entities.InformationDelivery) error {
	query := `
		INSERT INTO informationdeliveries (observations, signature_received, photo_delivery, reason_cancellation, idpackage)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query,
			info.Observation,
			info.SignatureReceived,
			info.PhotoDelivery,
			info.ReasonCancellation,
			info.PackageID,
		).Scan(&info.ID)
	} else {
		err = r.db.QueryRowContext(ctx, query,
			info.Observation,
			info.SignatureReceived,
			info.PhotoDelivery,
			info.ReasonCancellation,
			info.PackageID,
		).Scan(&info.ID)
	}

	if err != nil {
		log.Printf("ERROR creating information delivery for package %d: %v", info.PackageID, err)
		return fmt.Errorf("create information delivery: %w", err)
	}

	log.Printf("✓ Information delivery created: ID=%d, PackageID=%d", info.ID, info.PackageID)
	return nil
}

func (r *InformationDeliveryRepositoryPostgres) GetByID(ctx context.Context, id uint) (*entities.InformationDelivery, error) {
	query := `
		SELECT id, observations, signature_received, photo_delivery, reason_cancellation, idpackage
		FROM informationdeliveries
		WHERE id = $1
	`

	var info entities.InformationDelivery
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&info.ID,
		&info.Observation,
		&info.SignatureReceived,
		&info.PhotoDelivery,
		&info.ReasonCancellation,
		&info.PackageID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrInformationDeliveryNotFound
		}
		return nil, fmt.Errorf("get information delivery by id: %w", err)
	}

	return &info, nil
}

func (r *InformationDeliveryRepositoryPostgres) GetByPackageID(ctx context.Context, packageID uint) (*entities.InformationDelivery, error) {
	query := `
		SELECT id, observations, signature_received, photo_delivery, reason_cancellation, idpackage
		FROM informationdeliveries
		WHERE idpackage = $1
	`

	var info entities.InformationDelivery
	err := r.db.QueryRowContext(ctx, query, packageID).Scan(
		&info.ID,
		&info.Observation,
		&info.SignatureReceived,
		&info.PhotoDelivery,
		&info.ReasonCancellation,
		&info.PackageID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrInformationDeliveryNotFound
		}
		return nil, fmt.Errorf("get information delivery by package id: %w", err)
	}

	return &info, nil
}

func (r *InformationDeliveryRepositoryPostgres) Update(ctx context.Context, info *entities.InformationDelivery) error {
	query := `
		UPDATE informationdeliveries
		SET observations = $1, signature_received = $2, photo_delivery = $3, reason_cancellation = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		info.Observation,
		info.SignatureReceived,
		info.PhotoDelivery,
		info.ReasonCancellation,
		info.ID,
	)

	if err != nil {
		log.Printf("ERROR updating information delivery %d: %v", info.ID, err)
		return fmt.Errorf("update information delivery: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrInformationDeliveryNotFound
	}

	log.Printf("✓ Information delivery updated: ID=%d", info.ID)
	return nil
}
