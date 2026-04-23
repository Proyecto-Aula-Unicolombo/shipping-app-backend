package usepackages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

var ErrMissingRelatedInput = errors.New("missing related input")

// CreateOrFetchRelatedEntities crea (si ID==0) o verifica (si ID!=0) las entidades relacionadas dentro de la tx.
// Devuelve las entidades resultantes con sus IDs ya establecidos (creadas o existentes).
// Si alguna creación/consulta falla retorna error y el caller debe hacer rollback.
func CreateOrFetchRelatedEntitiesFromDTOs(
	ctx context.Context,
	tx *sql.Tx,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	input CreatePackageInput,
) (
	*entities.AddressPackage,
	*entities.ComercialInformation,
	*entities.Sender,
	*entities.Receiver,
	error,
) {
	// AddressPackage - Buscar por origin + destination
	if input.AddressPackage == nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: addresspackage", ErrMissingRelatedInput)
	}

	addrEntity, err := addressRepo.FindByRoute(ctx, input.AddressPackage.Origin, input.AddressPackage.Destination)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, nil, fmt.Errorf("find addresspackage: %w", err)
	}

	if addrEntity == nil {
		// No existe, crear nueva
		addrEntity = mapAddressInputToEntity(input.AddressPackage)
		if err := addressRepo.Create(ctx, tx, addrEntity); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("create addresspackage: %w", err)
		}
	}

	// ComercialInformation - Siempre crear nueva (específica del paquete)
	if input.ComercialInformation == nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: comercialinformation", ErrMissingRelatedInput)
	}

	cominfoEntity := mapComercialInputToEntity(input.ComercialInformation)
	if err := comercialRepo.Create(ctx, tx, cominfoEntity); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("create comercialinformation: %w", err)
	}

	senderEntity, err := senderRepo.FindByEmailOrDocument(ctx, "", input.SenderDocument)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, nil, fmt.Errorf("find sender: %w", err)
	}

	// Receiver - Buscar por email
	if input.Receiver == nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: receiver", ErrMissingRelatedInput)
	}

	receiverEntity, err := receiverRepo.FindByEmailWithTx(ctx, tx, input.Receiver.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, nil, fmt.Errorf("find receiver: %w", err)
	}

	if receiverEntity == nil {
		receiverEntity = mapReceiverInputToEntity(input.Receiver)
		if err := receiverRepo.Create(ctx, tx, receiverEntity); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("create receiver: %w", err)
		}
	}

	return addrEntity, cominfoEntity, senderEntity, receiverEntity, nil
}

func mapAddressInputToEntity(in *related.AdressPackageInput) *entities.AddressPackage {
	if in == nil {
		return nil
	}
	return &entities.AddressPackage{
		Origin:               in.Origin,
		Destination:          in.Destination,
		DeliveryInstructions: in.DeliveryInstructions,
	}
}

func mapComercialInputToEntity(in *related.ComercialInformationInput) *entities.ComercialInformation {
	if in == nil {
		return nil
	}
	return &entities.ComercialInformation{
		CostSending: in.CostSending,
		IsPaid:      in.IsPaid,
	}
}

func mapSenderInputToEntity(in *related.SenderInput) *entities.Sender {
	if in == nil {
		return nil
	}
	return &entities.Sender{
		Name:        in.Name,
		Document:    in.Document,
		Address:     in.Address,
		PhoneNumber: in.PhoneNumber,
		Email:       in.Email,
	}
}

func mapReceiverInputToEntity(in *related.ReceiverInput) *entities.Receiver {
	if in == nil {
		return nil
	}
	return &entities.Receiver{
		Name:        in.Name,
		LastName:    in.LastName,
		PhoneNumber: in.PhoneNumber,
		Email:       in.Email,
	}
}
