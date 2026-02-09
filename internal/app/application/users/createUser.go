package users

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/utils"
)

type CreateUserInput struct {
	Name        string
	LastName    string
	Email       string
	Password    string
	Role        string
	PhoneNumber string
	NumLicence  string
}

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidRole       = errors.New("invalid role")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type CreateUserUseCase struct {
	userRepo   repository.UserRepository
	driver     repository.DriverRepository
	tcProvider repository.TxProvider
}

func NewCreateUserUseCase(userRepo repository.UserRepository, driver repository.DriverRepository, tcProvider repository.TxProvider) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo, driver: driver, tcProvider: tcProvider}
}

func (us *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) error {
	if err := validateInput(input); err != nil {
		return err
	}
	userAlreadyExiste, _ := us.userRepo.GetUserByEmail(input.Email)

	if userAlreadyExiste != nil {
		return ErrUserAlreadyExists
	}

	passwordHashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return errors.New("error hashing password")
	}

	tx, err := us.tcProvider.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = us.tcProvider.RollbackTx(ctx, tx)
		}
	}()

	user := entities.User{
		Name:     input.Name,
		LastName: input.LastName,
		Email:    input.Email,
		Password: passwordHashed,
		Role:     input.Role,
	}

	if err := us.userRepo.CreateUserTx(tx, &user); err != nil {
		return err
	}

	if input.Role == "driver" {
		driver := entities.Driver{
			UserID:      user.ID,
			PhoneNumber: input.PhoneNumber,
			LicenseNo:   input.NumLicence,
		}
		if err := us.driver.CreateDriverTx(tx, &driver); err != nil {
			return errors.New("error creating driver")
		}
	}

	if err := us.tcProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true

	return nil
}

func validateInput(input CreateUserInput) error {
	if input.Name == "" || input.LastName == "" || input.Email == "" || input.Password == "" || input.Role == "" {
		return ErrInvalidInput
	}

	if len(input.Password) < 8 {
		return ErrPasswordTooShort
	}

	validRoles := map[string]bool{
		"coord":  true,
		"admin":  true,
		"driver": true,
	}

	if !validRoles[input.Role] {
		return ErrInvalidRole
	}

	return nil
}
