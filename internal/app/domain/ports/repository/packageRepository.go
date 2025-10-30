package repository

import (
	"context"
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type PackageRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pkg *entities.Package) error
}
