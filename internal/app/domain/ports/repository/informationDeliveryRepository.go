package repository

import (
	"context"
	"database/sql"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrInformationDeliveryNotFound = errors.New("information delivery not found")
)

type InformationDeliveryRepository interface {
	Create(ctx context.Context, tx *sql.Tx, info *entities.InformationDelivery) error
	GetByID(ctx context.Context, id uint) (*entities.InformationDelivery, error)
	GetByPackageID(ctx context.Context, packageID uint) (*entities.InformationDelivery, error)
	Update(ctx context.Context, info *entities.InformationDelivery) error
}
