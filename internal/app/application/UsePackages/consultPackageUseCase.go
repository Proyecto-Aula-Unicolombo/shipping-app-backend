package usepackages

import (
	"errors"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	"time"
)

type ResponsePackage struct {
	ID                   uint
	NumPackage           string
	Status               string
	DescriptionContent   *string
	Weight               *float64
	Dimension            *string
	DeclaredValue        *float64
	TypePackage          *string
	IsFragile            bool
	IdOrder              *uint
	CreatedAt            time.Time
	AddressPackage       *related.AdressPackageResponse
	ComercialInformation *related.ComercialInformationResponse
	Sender               *related.SenderResponse
	Receiver             *related.ReceiverResponse
	DeliveryInformation  *DeliveryInformationResponse
}

type DeliveryInformationResponse struct {
	ID                 uint
	Observation        *string
	SignatureReceived  *string
	PhotoDelivery      string
	ReasonCancellation *string
}
type ConsultPackageUseCase struct {
	packageRepo      repository.PackageRepository
	addressRepo      repository.AddressPackageRepository
	comercialRepo    repository.ComercialInformationRepository
	senderRepo       repository.SenderRepository
	receiverRepo     repository.ReceiverRepository
	infoDeliveryRepo repository.InformationDeliveryRepository
}

func NewConsultPackageUseCase(
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository,
	infoDeliveryRepo repository.InformationDeliveryRepository,
) *ConsultPackageUseCase {
	return &ConsultPackageUseCase{
		packageRepo:      packageRepo,
		addressRepo:      addressRepo,
		comercialRepo:    comercialRepo,
		senderRepo:       senderRepo,
		receiverRepo:     receiverRepo,
		infoDeliveryRepo: infoDeliveryRepo,
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

	addrEntity, cominfoEntity, senderEntity, receiverEntity, err := GetRelatedEntities(
		input.Ctx,
		uc.addressRepo,
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
		Status:             pkg.Status,
		DescriptionContent: pkg.DescriptionContent,
		Weight:             pkg.Weight,
		Dimension:          pkg.Dimension,
		DeclaredValue:      pkg.DeclaredValue,
		TypePackage:        pkg.TypePackage,
		IsFragile:          pkg.IsFragile,
		IdOrder:            pkg.OrderID,
		CreatedAt:          pkg.CreatedAt,
		AddressPackage: &related.AdressPackageResponse{
			Origin:               addrEntity.Origin,
			Destination:          addrEntity.Destination,
			DeliveryInstructions: addrEntity.DeliveryInstructions,
		},
		ComercialInformation: &related.ComercialInformationResponse{
			CostSending: cominfoEntity.CostSending,
			IsPaid:      cominfoEntity.IsPaid,
		},
		Sender: &related.SenderResponse{
			Name:  senderEntity.Name,
			Email: senderEntity.Email,
		},
		Receiver: &related.ReceiverResponse{
			Name:        receiverEntity.Name,
			LastName:    receiverEntity.LastName,
			PhoneNumber: receiverEntity.PhoneNumber,
			Email:       receiverEntity.Email,
		},
	}

	if input.PackageID != nil {
		infoDeliveryEntity, err := uc.infoDeliveryRepo.GetByPackageID(input.Ctx, *input.PackageID)
		if err != nil {
			if !errors.Is(err, repository.ErrInformationDeliveryNotFound) {
				return nil, err
			}
		} else {
			response.DeliveryInformation = &DeliveryInformationResponse{
				ID:                 infoDeliveryEntity.ID,
				Observation:        infoDeliveryEntity.Observation,
				SignatureReceived:  infoDeliveryEntity.SignatureReceived,
				PhotoDelivery:      infoDeliveryEntity.PhotoDelivery,
				ReasonCancellation: infoDeliveryEntity.ReasonCancellation,
			}
		}
	}

	return response, nil
}
