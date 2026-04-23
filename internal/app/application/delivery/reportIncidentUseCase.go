package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports"
	"shipping-app/internal/app/domain/ports/repository"
)

type ReportIncidentInput struct {
	PackageID          uint
	ReasonCancellation *string
	Observation        *string
	PhotoEvidence      *string
	Status             string
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
	infoDeliveryRepo  repository.InformationDeliveryRepository
	packageRepo       repository.PackageRepository
	txProvider        repository.TxProvider
	generateReportSvc ports.GenerateReportUseCasePort
}

func NewReportIncidentUseCase(
	infoDeliveryRepo repository.InformationDeliveryRepository,
	packageRepo repository.PackageRepository,
	txProvider repository.TxProvider,
	generateReportSvc ports.GenerateReportUseCasePort,

) *ReportIncidentUseCase {
	return &ReportIncidentUseCase{
		infoDeliveryRepo:  infoDeliveryRepo,
		packageRepo:       packageRepo,
		txProvider:        txProvider,
		generateReportSvc: generateReportSvc,
	}
}

func (uc *ReportIncidentUseCase) Execute(ctx context.Context, input ReportIncidentInput) (*ReportIncidentOutput, error) {
	// Validar input - solo campos obligatorios
	if input.PackageID == 0 {
		return nil, fmt.Errorf("%w: package_id is required", ErrInvalidIncidentInput)
	}

	if input.Status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidIncidentInput)
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
		existingInfo.ReasonCancellation = input.ReasonCancellation
		existingInfo.Observation = input.Observation
		if input.PhotoEvidence != nil {
			existingInfo.PhotoDelivery = *input.PhotoEvidence
		} else {
			existingInfo.PhotoDelivery = ""
		}

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
	photoDelivery := ""
	if input.PhotoEvidence != nil {
		photoDelivery = *input.PhotoEvidence
	}

	infoDelivery := &entities.InformationDelivery{
		ReasonCancellation: input.ReasonCancellation,
		Observation:        input.Observation,
		PhotoDelivery:      photoDelivery,
		PackageID:          input.PackageID,
	}

	if err := uc.infoDeliveryRepo.Create(ctx, tx, infoDelivery); err != nil {
		return nil, fmt.Errorf("create incident information: %w", err)
	}

	if err := uc.packageRepo.UpdatePackageStatusDelivery(ctx, tx, input.Status, input.PackageID); err != nil {
		return nil, fmt.Errorf("update package status delivery: %w", err)
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

	go func() {
		if err := uc.generateReportSvc.Execute(pkg.ID); err != nil {
			log.Printf("[USE CASE] error is generating  report pkg#%d: %v", pkg.ID, err)
		}
	}()

	return output, nil
}
