package repository

import (
	"context"
	"database/sql"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrPackageAlreadyExists = errors.New("package with this number already exists")
	ErrPackageNotFound      = errors.New("package not found")
)

type PackageRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pkg *entities.Package) error
	GetByNumPackage(ctx context.Context, numPackage int64) (*entities.Package, error)
	GetStatusPackageToCancel(ctx context.Context, id uint) (*entities.Package, error)
	DeletePackage(ctx context.Context, tx *sql.Tx, id uint) error
}
