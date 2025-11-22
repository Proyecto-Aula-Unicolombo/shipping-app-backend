package usepackages

import (
	"context"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListPackagesToCreateOrderInput struct {
	Ctx context.Context

	Limit  int
	Offset int
}

type ListPackagesToCreateOrderOutput struct {
	ID             uint
	NumPackage     string
	Status         string
	TypePackage    string
	AddressPackage *related.AdressPackageResponse
	Receiver       *related.ReceiverResponse
}

type ListPackagesToCreateOrderUseCase struct {
	packageRepo  repository.PackageRepository
	addressRepo  repository.AddressPackageRepository
	receiverRepo repository.ReceiverRepository
}

func NewListPackagesToCreateOrderUseCase(
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	receiverRepo repository.ReceiverRepository) *ListPackagesToCreateOrderUseCase {
	return &ListPackagesToCreateOrderUseCase{
		packageRepo:  packageRepo,
		addressRepo:  addressRepo,
		receiverRepo: receiverRepo,
	}
}

func (uc *ListPackagesToCreateOrderUseCase) Execute(input ListPackagesToCreateOrderInput) ([]*ListPackagesToCreateOrderOutput, int64, error) {
	total, err := uc.packageRepo.CountPackagesToCreateOrder(input.Ctx)
	if err != nil {
		return nil, 0, err
	}

	packagesList, err := uc.packageRepo.ListPackagestoCreateOrder(input.Ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, 0, err
	}

	var responsePackages []*ListPackagesToCreateOrderOutput
	for _, p := range packagesList {
		addrEntity, err := uc.addressRepo.GetByID(input.Ctx, p.AddressPackageID)
		if err != nil {
			return nil, 0, err
		}
		receiverEntity, err := uc.receiverRepo.GetByID(input.Ctx, p.ReceiverID)
		if err != nil {
			return nil, 0, err
		}
		responsePackage := &ListPackagesToCreateOrderOutput{
			ID:             p.ID,
			NumPackage:     p.NumPackage,
			Status:         p.Status,
			TypePackage:    *p.TypePackage,
			AddressPackage: &related.AdressPackageResponse{Origin: addrEntity.Origin, Destination: addrEntity.Destination},
			Receiver:       &related.ReceiverResponse{Name: receiverEntity.Name, LastName: receiverEntity.LastName, PhoneNumber: receiverEntity.PhoneNumber, Email: receiverEntity.Email},
		}
		responsePackages = append(responsePackages, responsePackage)
	}

	return responsePackages, total, nil
}
