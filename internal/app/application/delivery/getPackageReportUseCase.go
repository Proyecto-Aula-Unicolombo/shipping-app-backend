package delivery

import (
	"context"
	"fmt"

	"shipping-app/internal/app/domain/ports/repository"
)

type PackageReportResponse struct {
	ID                 uint    `json:"id"`
	PackageID          uint    `json:"package_id"`
	PackageNumPackage  string  `json:"package_num_package"`
	Observation        *string `json:"observation"`
	SignatureReceived  *string `json:"signature_received"`
	PhotoDelivery      string  `json:"photo_delivery"`
	ReasonCancellation *string `json:"reason_cancellation"`
	ReportType         string  `json:"report_type"` // "delivery" o "incident"
}

type GetPackageReportUseCase struct {
	infoDeliveryRepo repository.InformationDeliveryRepository
	packageRepo      repository.PackageRepository
}

func NewGetPackageReportUseCase(
	infoDeliveryRepo repository.InformationDeliveryRepository,
	packageRepo repository.PackageRepository,
) *GetPackageReportUseCase {
	return &GetPackageReportUseCase{
		infoDeliveryRepo: infoDeliveryRepo,
		packageRepo:      packageRepo,
	}
}

func (uc *GetPackageReportUseCase) ExecuteByPackageID(ctx context.Context, packageID uint) (*PackageReportResponse, error) {
	// Obtener información del paquete
	pkg, err := uc.packageRepo.GetByID(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}

	// Obtener información de entrega/incidente
	infoDelivery, err := uc.infoDeliveryRepo.GetByPackageID(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("get delivery information: %w", err)
	}

	// Determinar tipo de reporte
	reportType := "delivery"
	if infoDelivery.ReasonCancellation != nil && *infoDelivery.ReasonCancellation != "" {
		reportType = "incident"
	}

	response := &PackageReportResponse{
		ID:                 infoDelivery.ID,
		PackageID:          pkg.ID,
		PackageNumPackage:  pkg.NumPackage,
		Observation:        infoDelivery.Observation,
		SignatureReceived:  infoDelivery.SignatureReceived,
		PhotoDelivery:      infoDelivery.PhotoDelivery,
		ReasonCancellation: infoDelivery.ReasonCancellation,
		ReportType:         reportType,
	}

	return response, nil
}

func (uc *GetPackageReportUseCase) ExecuteByReportID(ctx context.Context, reportID uint) (*PackageReportResponse, error) {
	// Obtener información de entrega/incidente
	infoDelivery, err := uc.infoDeliveryRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("get delivery information: %w", err)
	}

	// Obtener información del paquete
	pkg, err := uc.packageRepo.GetByID(ctx, infoDelivery.PackageID)
	if err != nil {
		return nil, fmt.Errorf("get package: %w", err)
	}

	// Determinar tipo de reporte
	reportType := "delivery"
	if infoDelivery.ReasonCancellation != nil && *infoDelivery.ReasonCancellation != "" {
		reportType = "incident"
	}

	response := &PackageReportResponse{
		ID:                 infoDelivery.ID,
		PackageID:          pkg.ID,
		PackageNumPackage:  pkg.NumPackage,
		Observation:        infoDelivery.Observation,
		SignatureReceived:  infoDelivery.SignatureReceived,
		PhotoDelivery:      infoDelivery.PhotoDelivery,
		ReasonCancellation: infoDelivery.ReasonCancellation,
		ReportType:         reportType,
	}

	return response, nil
}
