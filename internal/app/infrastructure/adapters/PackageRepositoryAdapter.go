package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		rowErr = tx.QueryRowContext(ctx, query, args...).Scan(&pkg.ID, &pkg.CreatedAt)
	} else {
		rowErr = r.db.QueryRowContext(ctx, query, args...).Scan(&pkg.ID, &pkg.CreatedAt)
	}

	if rowErr != nil {
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
	}

	return nil
}

func (r *PackageRepositoryPostgres) GetByNumPackage(ctx context.Context, numPackage int64) (*entities.Package, error) {
	query := `
		SELECT id, numpackage
		FROM packages
		WHERE numpackage = $1
	`

	var pkg entities.Package
	var scanErr error

	scanErr = r.db.QueryRowContext(ctx, query, numPackage).Scan(
		&pkg.ID,
		&pkg.NumPackage,
	)

	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, repository.ErrPackageNotFound
		}
		return nil, fmt.Errorf("get package by numpackage: %w", scanErr)
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
