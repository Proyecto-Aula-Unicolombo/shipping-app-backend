package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"shipping-app/internal/app/domain/entities"
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

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = r.db.QueryRowContext(ctx, query, args...)
	}

	if err := row.Scan(&pkg.ID); err != nil {
		return fmt.Errorf("create package: %w", sanitizePGError(err))
	}
	return nil
}

func sanitizePGError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint") {
		return fmt.Errorf("unique_violation: %s", msg)
	}
	return err
}
