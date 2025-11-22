package orders

import (
	"context"
	"fmt"
	"time"

	usepackages "shipping-app/internal/app/application/UsePackages"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/app/application/vehicles"
	"shipping-app/internal/app/domain/ports/repository"
)

type DriverDTO struct {
	ID       uint
	Name     string
	LastName string
	Email    string
}
type ResponseOrder struct {
	ID          uint
	Status      string
	Observation *string
	TypeService string
	Driver      DriverDTO
	Vehicle     vehicles.VehiclesOutput
	CreateAt    time.Time
	AssignedAt  *time.Time
	Packages    []*usepackages.ResponsePackage
}

type GetOrderUseCase struct {
	orderRepo     repository.OrderRepository
	packageRepo   repository.PackageRepository
	addressRepo   repository.AddressPackageRepository
	comercialRepo repository.ComercialInformationRepository
	receiverRepo  repository.ReceiverRepository
	driverRepo    repository.DriverRepository
	vehicleRepo   repository.VehicleRepository
}

func NewGetOrderUseCase(
	orderRepo repository.OrderRepository,
	packageRepo repository.PackageRepository,
	addressRepo repository.AddressPackageRepository,
	comercialRepo repository.ComercialInformationRepository,
	receiverRepo repository.ReceiverRepository,
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepo:     orderRepo,
		packageRepo:   packageRepo,
		addressRepo:   addressRepo,
		comercialRepo: comercialRepo,
		receiverRepo:  receiverRepo,
		driverRepo:    driverRepo,
		vehicleRepo:   vehicleRepo,
	}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, id uint) (*ResponseOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	packages := []*usepackages.ResponsePackage{}
	for _, packageID := range order.PackageIDs {
		pkg, err := uc.packageRepo.GetByID(ctx, packageID)
		if err != nil {
			return nil, fmt.Errorf("get package: %w", err)
		}
		addrEntity, cominfoEntity, _, receiverEntity, err := usepackages.GetRelatedEntities(
			ctx,
			uc.addressRepo,
			uc.comercialRepo,
			nil,
			uc.receiverRepo,
			pkg,
		)
		if err != nil {
			return nil, usepackages.ErrGetRelatedEntities
		}

		response := &usepackages.ResponsePackage{
			ID:                 pkg.ID,
			NumPackage:         pkg.NumPackage,
			Status:             pkg.Status,
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
			ComercialInformation: &related.ComercialInformationResponse{
				CostSending: cominfoEntity.CostSending,
				IsPaid:      cominfoEntity.IsPaid,
			},
			Receiver: &related.ReceiverResponse{
				Name:        receiverEntity.Name,
				LastName:    receiverEntity.LastName,
				PhoneNumber: receiverEntity.PhoneNumber,
				Email:       receiverEntity.Email,
			},
		}

		packages = append(packages, response)
	}

	orderOutput := &ResponseOrder{
		ID:          order.ID,
		Status:      order.Status,
		Observation: order.Observation,
		TypeService: order.TypeService,
		CreateAt:    order.CreateAt,
		AssignedAt:  order.AssignedAt,
		Packages:    packages,
	}

	if order.DriverID != nil {
		driver, err := uc.driverRepo.GetByID(ctx, *order.DriverID)
		if err != nil {
			return nil, fmt.Errorf("get driver: %w", err)
		}

		orderOutput.Driver = DriverDTO{
			ID:       driver.ID,
			Name:     driver.User.Name,
			LastName: driver.User.LastName,
			Email:    driver.User.Email,
		}
	}

	if order.VehicleID != nil {
		vehicle, err := uc.vehicleRepo.GetByID(ctx, *order.VehicleID)
		if err != nil {
			return nil, fmt.Errorf("get vehicle: %w", err)
		}

		orderOutput.Vehicle = vehicles.VehiclesOutput{
			ID:          vehicle.ID,
			Plate:       vehicle.Plate,
			Brand:       vehicle.Brand,
			Model:       vehicle.Model,
			Color:       vehicle.Color,
			VehicleType: vehicle.VehicleType,
		}
	}

	return orderOutput, nil
}
