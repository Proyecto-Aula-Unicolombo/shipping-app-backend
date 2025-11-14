package users

import (
	"errors"
	"shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrUserNotFound = errors.New("Usuario no registrado")
	ErrInvalidID    = errors.New("ID inválido")
)

type UserOutput struct {
	ID       uint
	Name     string
	LastName string
	Email    string
	Role     string
	Driver   drivers.DriverDTO
}

type GetUser struct {
	repo       repository.UserRepository
	driverRepo repository.DriverRepository
}

func NewGetUser(repo repository.UserRepository, driverRepo repository.DriverRepository) *GetUser {
	return &GetUser{repo: repo, driverRepo: driverRepo}
}

func (uc *GetUser) Execute(id uint) (*UserOutput, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userOutput := &UserOutput{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Role:     user.Role,
	}

	if user.Role == "driver" {
		driver, err := uc.driverRepo.GetDriverByUserID(user.ID)
		if err != nil {
			return nil, err
		}
		userOutput.Driver = drivers.DriverDTO{
			ID:          driver.ID,
			PhoneNumber: driver.PhoneNumber,
			LicenseNo:   driver.LicenseNo,
		}
	}

	return userOutput, nil
}
