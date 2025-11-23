package repository

import (
	"context"
	"database/sql"

	"shipping-app/internal/app/domain/entities"
)

type AddressPackageRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.AddressPackage, error)
	FindByRoute(ctx context.Context, origin, destination string) (*entities.AddressPackage, error)
	Create(ctx context.Context, tx *sql.Tx, addr *entities.AddressPackage) error
}

type ComercialInformationRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.ComercialInformation, error)
	Create(ctx context.Context, tx *sql.Tx, info *entities.ComercialInformation) error
	Delete(ctx context.Context, tx *sql.Tx, id uint) error
}

type SenderRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Sender, error)
	FindByEmailOrDocument(ctx context.Context, email, document string) (*entities.Sender, error)
	Create(ctx context.Context, tx *sql.Tx, s *entities.Sender) error
}

type ReceiverRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Receiver, error)
	FindByEmail(ctx context.Context, email string) (*entities.Receiver, error)
	FindByEmailWithTx(ctx context.Context, tx *sql.Tx, email string) (*entities.Receiver, error)
	Create(ctx context.Context, tx *sql.Tx, r *entities.Receiver) error
}

type StatusDeliveryRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.StatusDelivery, error)
	Create(ctx context.Context, tx *sql.Tx, s *entities.StatusDelivery) error
	Delete(ctx context.Context, tx *sql.Tx, id uint) error
}
