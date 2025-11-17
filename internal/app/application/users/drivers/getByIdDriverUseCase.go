package drivers

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

type GetDriverOutput struct {
	ID          uint
	Name        string
	LastName    string
	Email       string
	PhoneNumber string
	NumLicence  string
	NumOrder    uint
	IsActive    bool
}

type GetDriversByIdUseCase struct {
	driverRepo repository.DriverRepository
}

var ErrNotFound = errors.New("driver not found")

func NewGetByIdDriverUseCase(driverRepo repository.DriverRepository) *GetDriversByIdUseCase {
	return &GetDriversByIdUseCase{
		driverRepo: driverRepo,
	}
}

func (uc *GetDriversByIdUseCase) Execute(ctx context.Context, id uint) (*GetDriverOutput, error) {
	driver, err := uc.driverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	driverOutput := GetDriverOutput{
		ID:          driver.ID,
		Name:        driver.User.Name,
		LastName:    driver.User.LastName,
		Email:       driver.User.Email,
		PhoneNumber: driver.PhoneNumber,
		NumLicence:  driver.LicenseNo,
		NumOrder:    driver.NumOrder,
		IsActive:    driver.IsActive,
	}

	return &driverOutput, nil

}
