package vehicles

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type CreateVehicleInput struct {
	Plate string
	Brand string
	Model string
	Color string
	VehicleType string

}

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
)

type CreateVehicleUseCase struct {
	vehicleRepo repository.VehicleRepository
	txProvider  repository.TxProvider
}

func NewCreateVehicleUseCase(
	vehicleRepo repository.VehicleRepository,
	txProvider repository.TxProvider,
) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{
		vehicleRepo: vehicleRepo,
		txProvider:  txProvider,
	}
}

func (uc *CreateVehicleUseCase) Execute(ctx context.Context, input CreateVehicleInput) error {
	if err := validateVehicleInput(input); err != nil {
		return err
	}

	tx, err := uc.txProvider.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = uc.txProvider.RollbackTx(ctx, tx)
		}
	}()

	vehicle := entities.Vehicle{
		Plate: input.Plate,
		Brand: input.Brand,
		Model: input.Model,
		Color: input.Color,
		VehicleType: input.VehicleType,

	}

	if err := uc.vehicleRepo.CreateVehicleTx(tx, &vehicle); err != nil {
		return err
	}

	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func validateVehicleInput(input CreateVehicleInput) error {
	if input.Plate == "" || input.Brand == "" || input.Model == "" || input.Color == "" || input.VehicleType == "" {
		return ErrInvalidInput
	}
	return nil
}
