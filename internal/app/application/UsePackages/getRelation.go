package usepackages

import (
	"context"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

func GetRelatedEntities(
	ctx context.Context,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	packageEntity *entities.Package,
) (
	*entities.AddressPackage,
	*entities.ComercialInformation,
	*entities.Sender,
	*entities.Receiver,
	error,
) {
	addrEntity, err := addressRepo.GetByID(ctx, packageEntity.AddressPackageID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	cominfoEntity, err := comercialRepo.GetByID(ctx, packageEntity.ComercialInformationID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// puede ser necesario o no, solo es por prueba, ya que el que deberia recibir es el mismo sender.
	senderEntity, err := senderRepo.GetByID(ctx, packageEntity.SenderID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	receiverEntity, err := receiverRepo.GetByID(ctx, packageEntity.ReceiverID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return addrEntity, cominfoEntity, senderEntity, receiverEntity, nil
}
