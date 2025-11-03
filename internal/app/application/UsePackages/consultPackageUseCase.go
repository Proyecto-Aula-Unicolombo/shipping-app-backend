package usepackages

import (
	"context"
	"errors"
	"shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ConsultPackageInput struct {
	CTX context.Context

	PackageID  *uint  `json:"package_id,omitempty"`  // Para UI (interno)
	NumPackage *int64 `json:"num_package,omitempty"` // Para API Key (externo)

	AuthType string `json:"-"` // "jwt" o "api_key"
	UserRole string `json:"-"` // "coordinator", "driver"
	DriverID *uint  `json:"-"` // ID del conductor (para filtrar)
	SenderID *uint  `json:"-"` // ID del sender (para filtrar)
}

type ResponsePackage struct {
	ID                   uint
	NumPackage           int64
	StartStatus          string
	DescriptionContent   *string
	Weight               *float64
	Dimension            *float64
	DeclaredValue        *float64
	TypePackage          *string
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

func (uc *ConsultPackageUseCase) Execute(input ConsultPackageInput) (*ResponsePackage, error) {

	if input.PackageID == nil && input.NumPackage == nil {
		return nil, ErrInvalidSearchCriteria
	}

	var pkg *entities.Package
	var err error

	// Buscar paquete según el criterio
	if input.PackageID != nil {
		// Búsqueda por ID (típicamente desde UI)
		pkg, err = uc.packageRepo.GetByID(input.CTX, *input.PackageID)
	} else {
		// Búsqueda por NumPackage (típicamente desde API Key)
		pkg, err = uc.packageRepo.GetByNumPackage(input.CTX, *input.NumPackage)
	}

	if err != nil {
		return nil, err
	}

	// Verificar permisos de acceso
	if err := uc.checkAccess(pkg, input); err != nil {
		return nil, err
	}

	addrEntity, statusEntity, cominfoEntity, senderEntity, receiverEntity, err := GetRelatedEntities(
		input.CTX,
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

// checkAccess verifica si el usuario/sender tiene acceso al paquete
func (uc *ConsultPackageUseCase) checkAccess(pkg *entities.Package, input ConsultPackageInput) error {
	switch input.AuthType {
	case "api_key":
		if input.SenderID != nil && pkg.SenderID != *input.SenderID {
			return ErrAccessDenied
		}

	case "jwt":
		switch input.UserRole {
		case "coordinator":
			return nil

		case "driver":
			return nil
		default:
			return ErrAccessDenied
		}
	}

	return nil
}
