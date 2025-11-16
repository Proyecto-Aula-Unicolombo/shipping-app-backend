package usepackages

import (
	"errors"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ResponsePackage struct {
	ID                   uint
	NumPackage           string
	StartStatus          string
	DescriptionContent   *string
	Weight               *float64
	Dimension            *string
	DeclaredValue        *float64
	TypePackage          *string
	IsFragile            bool
	AddressPackage       *related.AdressPackageResponse
	StatusDelivery       *related.StatusDeliveryResponse
	ComercialInformation *related.ComercialInformationResponse
	Sender               *related.SenderResponse
	Receiver             *related.ReceiverResponse
}
type ConsultPackageUseCase struct {
	packageRepo   repository.PackageRepository
	addressRepo   repository.AddressPackageRepository
	comercialRepo repository.ComercialInformationRepository
	senderRepo    repository.SenderRepository
	receiverRepo  repository.ReceiverRepository
	statusRepo    repository.StatusDeliveryRepository
}

func NewConsultPackageUseCase(
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	statusRepo repository.StatusDeliveryRepository,
) *ConsultPackageUseCase {
	return &ConsultPackageUseCase{
		packageRepo:   packageRepo,
		addressRepo:   addressRepo,
		comercialRepo: comercialRepo,
		senderRepo:    senderRepo,
		receiverRepo:  receiverRepo,
		statusRepo:    statusRepo,
	}
}

var (
	ErrGetRelatedEntities    = errors.New("error getting related entities")
	ErrInvalidSearchCriteria = errors.New("must provide either package_id or num_package")
	ErrAccessDenied          = errors.New("access denied to this package")
)

func (uc *ConsultPackageUseCase) Execute(input CheckAccessInput) (*ResponsePackage, error) {

	if input.PackageID == nil && input.NumPackage == nil {
		return nil, ErrInvalidSearchCriteria
	}

	var pkg *entities.Package
	var err error

	// Buscar paquete según el criterio
	if input.PackageID != nil {
		// Búsqueda por ID (típicamente desde UI)
		pkg, err = uc.packageRepo.GetByID(input.Ctx, *input.PackageID)
	} else {
		// Búsqueda por NumPackage (típicamente desde API Key)
		pkg, err = uc.packageRepo.GetByNumPackage(input.Ctx, *input.NumPackage)
	}

	if err != nil {
		return nil, err
	}

	// Verificar permisos de acceso
	if err := CheckAccess(pkg, input); err != nil {
		return nil, err
	}

	addrEntity, statusEntity, cominfoEntity, senderEntity, receiverEntity, err := GetRelatedEntities(
		input.Ctx,
		uc.addressRepo,
		uc.statusRepo,
		uc.comercialRepo,
		uc.senderRepo,
		uc.receiverRepo,
		pkg,
	)
	if err != nil {
		return nil, ErrGetRelatedEntities
	}

	// Construir respuesta
	response := &ResponsePackage{
		ID:                 pkg.ID,
		NumPackage:         pkg.NumPackage,
		StartStatus:        pkg.StartStatus,
		DescriptionContent: pkg.DescriptionContent,
		Weight:             pkg.Weight,
		Dimension:          pkg.Dimension,
		DeclaredValue:      pkg.DeclaredValue,
		TypePackage:        pkg.TypePackage,
		IsFragile:          pkg.IsFragile,
		AddressPackage: &related.AdressPackageResponse{
			Origin:               addrEntity.Origin,
			Destination:          addrEntity.Destination,
			DeliveryInstructions: addrEntity.DeliveryInstructions,
		},
		StatusDelivery: &related.StatusDeliveryResponse{
			Status:                statusEntity.Status,
			Priority:              statusEntity.Priority,
			DateEstimatedDelivery: statusEntity.DateEstimatedDelivery,
		},
		ComercialInformation: &related.ComercialInformationResponse{
			CostSending: cominfoEntity.CostSending,
			IsPaid:      cominfoEntity.IsPaid,
		},
		Sender: &related.SenderResponse{
			Name:        senderEntity.Name,
			Document:    senderEntity.Document,
			Address:     senderEntity.Address,
			PhoneNumber: senderEntity.PhoneNumber,
			Email:       senderEntity.Email,
		},
		Receiver: &related.ReceiverResponse{
			Name:        receiverEntity.Name,
			LastName:    receiverEntity.LastName,
			PhoneNumber: receiverEntity.PhoneNumber,
			Email:       receiverEntity.Email,
		},
	}

	return response, nil
}
