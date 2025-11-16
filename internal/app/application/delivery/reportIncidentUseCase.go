package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ReportIncidentInput struct {
	PackageID          uint
	ReasonCancellation string
	Observation        *string
	PhotoEvidence      string
}

type ReportIncidentOutput struct {
	ID        uint
	PackageID uint
	CreatedAt time.Time
}

var (
	ErrInvalidIncidentInput = errors.New("invalid incident input")
	ErrPackageHasIncident   = errors.New("package already has incident report")
)

type ReportIncidentUseCase struct {
	infoDeliveryRepo repository.InformationDeliveryRepository
	packageRepo      repository.PackageRepository
	txProvider       repository.TxProvider
}

func NewReportIncidentUseCase(
	infoDeliveryRepo repository.InformationDeliveryRepository,
	packageRepo repository.PackageRepository,
	txProvider repository.TxProvider,
) *ReportIncidentUseCase {
	return &ReportIncidentUseCase{
		infoDeliveryRepo: infoDeliveryRepo,
		packageRepo:      packageRepo,
		txProvider:       txProvider,
	}
}

func (uc *ReportIncidentUseCase) Execute(ctx context.Context, input ReportIncidentInput) (*ReportIncidentOutput, error) {
	// Validar input
	if input.PackageID == 0 || input.ReasonCancellation == "" {
		return nil, ErrInvalidIncidentInput
	}

	if input.PhotoEvidence == "" {
		return nil, fmt.Errorf("%w: photo evidence is required", ErrInvalidIncidentInput)
	}

	// Verificar que el paquete existe
	pkg, err := uc.packageRepo.GetByID(ctx, input.PackageID)
	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}

	// Verificar que no exista reporte de incidente previo
	existingInfo, _ := uc.infoDeliveryRepo.GetByPackageID(ctx, input.PackageID)
	if existingInfo != nil && existingInfo.ReasonCancellation != nil {
		return nil, ErrPackageHasIncident
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

	// Crear o actualizar información de incidente
	if existingInfo != nil {
		// Actualizar existente
		existingInfo.ReasonCancellation = &input.ReasonCancellation
		existingInfo.Observation = input.Observation
		existingInfo.PhotoDelivery = input.PhotoEvidence

		if err := uc.infoDeliveryRepo.Update(ctx, existingInfo); err != nil {
			return nil, fmt.Errorf("update incident information: %w", err)
		}

		if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
			return nil, fmt.Errorf("commit tx: %w", err)
		}
		committed = true

		return &ReportIncidentOutput{
			ID:        existingInfo.ID,
			PackageID: pkg.ID,
			CreatedAt: time.Now(),
		}, nil
	}

	// Crear nueva información de incidente
	infoDelivery := &entities.InformationDelivery{
		ReasonCancellation: &input.ReasonCancellation,
		Observation:        input.Observation,
		PhotoDelivery:      input.PhotoEvidence,
		PackageID:          input.PackageID,
	}

	if err := uc.infoDeliveryRepo.Create(ctx, tx, infoDelivery); err != nil {
		return nil, fmt.Errorf("create incident information: %w", err)
	}

	// Commit transacción
	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	output := &ReportIncidentOutput{
		ID:        infoDelivery.ID,
		PackageID: pkg.ID,
		CreatedAt: time.Now(),
	}

	return output, nil
}
