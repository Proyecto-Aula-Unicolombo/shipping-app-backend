package usepackages

import (
	"context"
	"errors"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListPackagesInput struct {
	Ctx context.Context

	Limit  int
	Offset int

	AuthType string
	SenderID *uint
	UserRole string
}

type ListPackagesUseCase struct {
	packageRepo   repository.PackageRepository
	addressRepo   repository.AddressPackageRepository
	comercialRepo repository.ComercialInformationRepository
	senderRepo    repository.SenderRepository
	receiverRepo  repository.ReceiverRepository
}

func NewListPackagesUseCase(
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	senderRepo repository.SenderRepository,
	receiverRepo repository.ReceiverRepository) *ListPackagesUseCase {
	return &ListPackagesUseCase{
		packageRepo:   packageRepo,
		addressRepo:   addressRepo,
		comercialRepo: comercialRepo,
		senderRepo:    senderRepo,
		receiverRepo:  receiverRepo,
	}
}

var ErrNoPackagesFound = errors.New("no packages found")

func (uc *ListPackagesUseCase) Execute(input ListPackagesInput) ([]*ResponsePackage, int64, error) {
	var total int64
	var pkg []*entities.Package
	var err error

	if input.SenderID == nil && input.AuthType == "jwt" {
		pkg, err = uc.packageRepo.ListPackages(input.Ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, 0, err
		}
		if len(pkg) == 0 {
			return nil, 0, ErrNoPackagesFound
		}
		total = int64(len(pkg))
	} else {
		pkg, err = uc.packageRepo.ListPackagesBySenderID(input.Ctx, *input.SenderID, input.Limit, input.Offset)
		if err != nil {
			return nil, 0, err
		}
		if len(pkg) == 0 {
			return nil, 0, ErrNoPackagesFound
		}
		total = int64(len(pkg))

	}

	allowedPkg, err := CheckBulkAccess(pkg, CheckAccessInput{Ctx: input.Ctx, AuthType: input.AuthType, SenderID: input.SenderID, UserRole: input.UserRole})
	if err != nil {
		return nil, 0, err
	}

	var responsePackages []*ResponsePackage
	for _, p := range allowedPkg {
		addrEntity, cominfoEntity, senderEntity, receiverEntity, err := GetRelatedEntities(
			input.Ctx,
			uc.addressRepo,
			uc.comercialRepo,
			uc.senderRepo,
			uc.receiverRepo,
			p,
		)
		if err != nil {
			return nil, 0, ErrGetRelatedEntities
		}
		responsePackage := &ResponsePackage{
			ID:                 p.ID,
			NumPackage:         p.NumPackage,
			Status:             p.Status,
			DescriptionContent: p.DescriptionContent,
			Weight:             p.Weight,
			Dimension:          p.Dimension,
			DeclaredValue:      p.DeclaredValue,
			TypePackage:        p.TypePackage,
			IsFragile:          p.IsFragile,
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
		responsePackages = append(responsePackages, responsePackage)

	}

	return responsePackages, total, nil
}
