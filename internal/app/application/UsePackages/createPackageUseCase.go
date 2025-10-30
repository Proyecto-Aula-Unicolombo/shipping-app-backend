package usepackages

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	services "shipping-app/internal/app/domain/services/package"
	"time"
)

var (
	ErrInvalidInput               = errors.New("invalid input")
	ErrRelatedEntityMustProvideID = errors.New("related entities must provide non-zero ID")
	ErrRelatedEntityNotFound      = errors.New("related entity not found")
	ErrBusinessRuleViolation      = errors.New("business rule violation")
)

type CreatePackageInput struct {
	NumPackage           int64
	StartStatus          string
	DescriptionContent   *string
	Weight               *float64
	Dimension            *float64
	DeclaredValue        *float64
	TypePackage          *string
	IsFragile            bool
	AddressPackage       *related.AdressPackageInput
	StatusDelivery       *related.StatusDeliveryInput
	ComercialInformation *related.ComercialInformationInput
	Sender               *related.SenderInput
	Receiver             *related.ReceiverInput
}

type CreatePackageOutput struct {
	ID         uint
	NumPackage int64
	CreatedAt  time.Time
}

type CreatePackageUseCase struct {
	txProvider    repository.TxProvider
	packageRepo   repository.PackageRepository
	addressRepo   repository.AddressPackageRepository
	comercialRepo repository.ComercialInformationRepository
	senderRepo    repository.SenderRepository
	receiverRepo  repository.ReceiverRepository
	statusRepo    repository.StatusDeliveryRepository
	domainSvc     *services.ValidatePackageService
}

func NewCreatePackageUseCase(
	txProvider repository.TxProvider,
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	statusRepo repository.StatusDeliveryRepository,
	domainSvc *services.ValidatePackageService,
) *CreatePackageUseCase {
	return &CreatePackageUseCase{
		txProvider:    txProvider,
		packageRepo:   packageRepo,
		addressRepo:   addressRepo,
		comercialRepo: comercialRepo,
		senderRepo:    senderRepo,
		receiverRepo:  receiverRepo,
		statusRepo:    statusRepo,
		domainSvc:     domainSvc,
	}
}

func (uc *CreatePackageUseCase) Execute(ctx context.Context, input CreatePackageInput) (*CreatePackageOutput, error) {

	if err := ValidateCreateInput(input); err != nil {
		return nil, err
	}
	if err := ValidateBusinessRules(uc.domainSvc, input); err != nil {
		return nil, err
	}

	//  Iniciar transacción
	tx, err := uc.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = uc.txProvider.RollbackTx(ctx, tx)
		}
	}()

	// Crear o verificar relacionadas dentro de la tx
	addr, status, cominfo, sender, receiver, err := CreateOrFetchRelatedEntitiesFromDTOs(
		ctx,
		tx,
		uc.addressRepo,
		uc.statusRepo,
		uc.comercialRepo,
		uc.senderRepo,
		uc.receiverRepo,
		input,
	)
	if err != nil {
		return nil, err
	}
	//  Preparar entidad Package y persistir
	pkg := &entities.Package{
		NumPackage:             input.NumPackage,
		StartStatus:            input.StartStatus,
		DescriptionContent:     input.DescriptionContent,
		Weight:                 input.Weight,
		Dimension:              input.Dimension,
		DeclaredValue:          input.DeclaredValue,
		TypePackage:            input.TypePackage,
		IsFragile:              input.IsFragile,
		AddressPackageID:       addr.ID,
		StatusDeliveryID:       status.ID,
		ComercialInformationID: cominfo.ID,
		SenderID:               sender.ID,
		ReceiverID:             receiver.ID,
	}

	if err := uc.packageRepo.Create(ctx, tx, pkg); err != nil {
		return nil, fmt.Errorf("create package: %w", err)
	}

	//  Commit
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return &CreatePackageOutput{
		ID:         pkg.ID,
		NumPackage: pkg.NumPackage,
	}, nil
}
