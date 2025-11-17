package drivers

import "shipping-app/internal/app/domain/ports/repository"

type UpdateStatusDriverUseCase struct {
	driverRepo repository.DriverRepository
}

func NewUpdateStatusDriverUseCase(driverRepo repository.DriverRepository) *UpdateStatusDriverUseCase {
	return &UpdateStatusDriverUseCase{
		driverRepo: driverRepo,
	}
}

func (uc *UpdateStatusDriverUseCase) Execute(driverID uint, isActive bool) error {
	if err := uc.driverRepo.UpdateDriverStatus(driverID, isActive); err != nil {
		return err
	}
	return nil
}
