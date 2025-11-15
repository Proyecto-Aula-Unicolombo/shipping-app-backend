package users

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/utils"
)

type UpdateUserUseCase struct {
	repo       repository.UserRepository
	driverRepo repository.DriverRepository
	txProvider repository.TxProvider
}

func NewUpdateUserUseCase(repo repository.UserRepository, driverRepo repository.DriverRepository, txProvider repository.TxProvider) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo, driverRepo: driverRepo, txProvider: txProvider}
}

type UpdateUserInput struct {
	ID       uint
	Name     string
	LastName string
	Email    string
	Password string
	Role     string
	Driver   drivers.DriverUpdateDTO
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) error {
	if input.ID == 0 {
		return errors.New("ID inválido")
	}

	existingUser, err := uc.repo.GetUserByID(input.ID)
	if err != nil {
		return errors.New("usuario no encontrado")
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

	if existingUser.Role == "driver" {
		driver, err := uc.driverRepo.GetDriverByUserID(input.ID)
		if err != nil {
			return errors.New("driver no found")
		}

		if input.Driver.PhoneNumber != "" {
			driver.PhoneNumber = input.Driver.PhoneNumber
		}
		if input.Driver.LicenseNo != "" {
			driver.LicenseNo = input.Driver.LicenseNo
		}

		driver.UserID = input.ID
		err = uc.driverRepo.UpdateDriverTx(tx, driver)
		if err != nil {
			return errors.New("error updating driver")
		}
	}

	if input.Name != "" {
		existingUser.Name = input.Name
	}
	if input.LastName != "" {
		existingUser.LastName = input.LastName
	}
	if input.Email != "" {
		existingUser.Email = input.Email
	}
	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		existingUser.Password = hashedPassword
	}
	if input.Role != "" {
		existingUser.Role = input.Role
	}

	err = uc.repo.UpdateUser(tx, existingUser)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}
