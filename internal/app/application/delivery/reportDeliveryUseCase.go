package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ReportDeliveryInput struct {
	PackageID         uint
	Observation       *string
	SignatureReceived *string
	PhotoDelivery     string
}

type ReportDeliveryOutput struct {
	ID        uint
	PackageID uint
	CreatedAt time.Time
}

var (
	ErrPackageAlreadyDelivered = errors.New("package already has delivery information")
	ErrInvalidDeliveryInput    = errors.New("invalid delivery input")
)

type ReportDeliveryUseCase struct {
	infoDeliveryRepo repository.InformationDeliveryRepository
	packageRepo      repository.PackageRepository
	txProvider       repository.TxProvider
}

func NewReportDeliveryUseCase(
	infoDeliveryRepo repository.InformationDeliveryRepository,
	packageRepo repository.PackageRepository,
	txProvider repository.TxProvider,
) *ReportDeliveryUseCase {
	return &ReportDeliveryUseCase{
		infoDeliveryRepo: infoDeliveryRepo,
		packageRepo:      packageRepo,
		txProvider:       txProvider,
	}
}

func (uc *ReportDeliveryUseCase) Execute(ctx context.Context, input ReportDeliveryInput) (*ReportDeliveryOutput, error) {
	// Validar input
	if input.PackageID == 0 {
		return nil, ErrInvalidDeliveryInput
	}

	if input.PhotoDelivery == "" {
		return nil, fmt.Errorf("%w: photo delivery is required", ErrInvalidDeliveryInput)
	}

	// Verificar que el paquete existe
	pkg, err := uc.packageRepo.GetByID(ctx, input.PackageID)
	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}

	// Verificar que no exista información de entrega previa
	existingInfo, _ := uc.infoDeliveryRepo.GetByPackageID(ctx, input.PackageID)
	if existingInfo != nil {
		return nil, ErrPackageAlreadyDelivered
	}

	// Iniciar transacción
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

	// Crear información de entrega
	infoDelivery := &entities.InformationDelivery{
		Observation:       input.Observation,
		SignatureReceived: input.SignatureReceived,
		PhotoDelivery:     input.PhotoDelivery,
		PackageID:         input.PackageID,
	}

	if err := uc.infoDeliveryRepo.Create(ctx, tx, infoDelivery); err != nil {
		return nil, fmt.Errorf("create delivery information: %w", err)
	}

	if err := uc.packageRepo.UpdatePackageStatusDelivery(ctx, tx, "entregado", input.PackageID); err != nil {
		return nil, fmt.Errorf("update package status delivery: %w", err)
	}

	// Commit transacción
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	output := &ReportDeliveryOutput{
		ID:        infoDelivery.ID,
		PackageID: pkg.ID,
		CreatedAt: time.Now(),
	}

	return output, nil
}
