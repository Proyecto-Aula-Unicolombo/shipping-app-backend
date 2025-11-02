package usepackages

import (
	"context"
	"database/sql"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

func GetRelatedEntities(
	ctx context.Context,
	tx *sql.Tx,
	addressRepo repository.AddressPackageRepository,
	statusRepo repository.StatusDeliveryRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	packageEntity *entities.Package,
) (
	*entities.AddressPackage,
	*entities.StatusDelivery,
	*entities.ComercialInformation,
	*entities.Sender,
	*entities.Receiver,
	error,
) {
	addrEntity, err := addressRepo.GetByID(ctx, tx, packageEntity.AddressPackageID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	statusEntity, err := statusRepo.GetByID(ctx, tx, packageEntity.StatusDeliveryID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	cominfoEntity, err := comercialRepo.GetByID(ctx, tx, packageEntity.ComercialInformationID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// puede ser necesario o no, solo es por prueba, ya que el que deberia recibir es el mismo sender.
	senderEntity, err := senderRepo.GetByID(ctx, tx, packageEntity.SenderID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	receiverEntity, err := receiverRepo.GetByID(ctx, tx, packageEntity.ReceiverID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return addrEntity, statusEntity, cominfoEntity, senderEntity, receiverEntity, nil
}
