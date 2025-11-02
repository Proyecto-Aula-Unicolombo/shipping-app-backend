package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/lib/pq"
)

type PackageRepositoryPostgres struct {
	db *sql.DB
}

func NewPackageRepositoryPostgres(db *sql.DB) *PackageRepositoryPostgres {
	return &PackageRepositoryPostgres{db: db}
}

func (r *PackageRepositoryPostgres) Create(ctx context.Context, tx *sql.Tx, pkg *entities.Package) error {
	query := `
		INSERT INTO packages (
			numpackage,
			startstatus,
			descriptioncontent,
			weight,
			dimension,
			declared_value,
			type_package,
			is_fragile,
			idaddresspackage,
			idstatusdelivery,
			idcomercialinformation,
			idsender,
			idreceivers
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id
	`

	args := []interface{}{
		pkg.NumPackage,
		pkg.StartStatus,
		pkg.DescriptionContent,
		pkg.Weight,
		pkg.Dimension,
		pkg.DeclaredValue,
		pkg.TypePackage,
		pkg.IsFragile,
		pkg.AddressPackageID,
		pkg.StatusDeliveryID,
		pkg.ComercialInformationID,
		pkg.SenderID,
		pkg.ReceiverID,
	}

	var rowErr error
	if tx != nil {
		rowErr = tx.QueryRowContext(ctx, query, args...).Scan(&pkg.ID)
	} else {
		rowErr = r.db.QueryRowContext(ctx, query, args...).Scan(&pkg.ID)
	}

	if rowErr == nil {
		log.Printf("✓ Package created successfully: ID=%d, NumPackage=%d", pkg.ID, pkg.NumPackage)
		return nil
	}

	if isDuplicatePackageError(rowErr) {

		existingPkg, err := r.GetByNumPackage(ctx, pkg.NumPackage)
		if err != nil && existingPkg == nil {
			return &PackageConflictError{
				NumPackage: pkg.NumPackage,
				ExistingID: 0,
			}
		}
		return &PackageConflictError{
			NumPackage: pkg.NumPackage,
			ExistingID: existingPkg.ID,
		}
	}

	return fmt.Errorf("package create: %w", rowErr)
}

func (r *PackageRepositoryPostgres) GetByNumPackage(ctx context.Context, numPackage int64) (*entities.Package, error) {
	query := `
		SELECT id, numpackage, startstatus, descriptioncontent, weight, dimension, declared_value, type_package, is_fragile,
		       idaddresspackage, idstatusdelivery, idcomercialinformation, idsender, idreceivers, created_at, updated_at
		FROM packages
		WHERE numpackage = $1
	`

	var pkg entities.Package
	var scanErr error

	scanErr = r.db.QueryRowContext(ctx, query, numPackage).Scan(
		&pkg.ID,
		&pkg.NumPackage,
		&pkg.StartStatus,
		&pkg.DescriptionContent,
		&pkg.Weight,
		&pkg.Dimension,
		&pkg.DeclaredValue,
		&pkg.TypePackage,
		&pkg.IsFragile,
		&pkg.AddressPackageID,
		&pkg.StatusDeliveryID,
		&pkg.ComercialInformationID,
		&pkg.SenderID,
		&pkg.ReceiverID,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	)

	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, repository.ErrPackageNotFound
		}
		return nil, fmt.Errorf("get package by numpackage: %w", scanErr)
	}
	return &pkg, nil
}

func (r *PackageRepositoryPostgres) GetStatusPackageToCancel(ctx context.Context, id uint) (*entities.Package, error) {
	query := `
		SELECT startstatus
		FROM packages
		WHERE id = $1
	`

	var pkg entities.Package
	var scanErr error
	scanErr = r.db.QueryRowContext(ctx, query, id).Scan(
		&pkg.StartStatus,
	)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, repository.ErrPackageNotFound
		}
		return nil, fmt.Errorf("get package status to cancel: %w", scanErr)
	}
	return &pkg, nil
}

func (r *PackageRepositoryPostgres) DeletePackage(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `DELETE FROM packages WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, id)
	}
	return err
}

func (r *PackageRepositoryPostgres) GetByID(ctx context.Context, tx *sql.Tx, id uint) (*entities.Package, error) {
	query := `
		SELECT id, numpackage, startstatus, descriptioncontent, weight, dimension, declared_value, type_package, is_fragile,
		       idaddresspackage, idstatusdelivery, idcomercialinformation, idsender, idreceivers, created_at, updated_at
		FROM packages
		WHERE id = $1
	`
	var pkg entities.Package
	var scanErr error
	if tx != nil {
		scanErr = tx.QueryRowContext(ctx, query, id).Scan(
			&pkg.ID,
			&pkg.NumPackage,
			&pkg.StartStatus,
			&pkg.DescriptionContent,
			&pkg.Weight,
			&pkg.Dimension,
			&pkg.DeclaredValue,
			&pkg.TypePackage,
			&pkg.IsFragile,
			&pkg.AddressPackageID,
			&pkg.StatusDeliveryID,
			&pkg.ComercialInformationID,
			&pkg.SenderID,
			&pkg.ReceiverID,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		)
	}

	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, repository.ErrPackageNotFound
		}
		return nil, fmt.Errorf("get package by id: %w", scanErr)
	}

	return &pkg, nil
}

type PackageConflictError struct {
	NumPackage int64
	ExistingID uint
}

func (e *PackageConflictError) Error() string {
	if e.ExistingID == 0 {
		return fmt.Sprintf("package with numpackage %d already exists (ID unknown)", e.NumPackage)
	}
	return fmt.Sprintf("package with numpackage %d already exists with id %d", e.NumPackage, e.ExistingID)
}

func isDuplicatePackageError(err error) bool {
	if err == nil {
		return false
	}

	// PostgreSQL unique violation error code: 23505
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		isDuplicate := pqErr.Code == "23505"
		return isDuplicate && (pqErr.Constraint == "packages_numpackage_key" ||
			strings.Contains(pqErr.Message, "packages_numpackage_key"))
	}

	// Fallback: buscar en el mensaje de error
	errMsg := err.Error()
	isDuplicate := strings.Contains(errMsg, "duplicate key") &&
		strings.Contains(errMsg, "packages_numpackage_key")

	return isDuplicate
}
