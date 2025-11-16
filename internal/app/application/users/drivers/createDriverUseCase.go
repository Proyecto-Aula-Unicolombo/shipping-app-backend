package drivers

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/utils"
)

type CreateDriverInput struct {
	Name       string
	LastName   string
	Email      string
	Phone      string
	NumLicence string
}

type CreateDriverUseCase struct {
	userRepo   repository.UserRepository
	driver     repository.DriverRepository
	txProvider repository.TxProvider
}

func NewCreateDriverUseCase(userRepo repository.UserRepository, driver repository.DriverRepository, txProvider repository.TxProvider) *CreateDriverUseCase {
	return &CreateDriverUseCase{
		userRepo:   userRepo,
		driver:     driver,
		txProvider: txProvider,
	}
}

var ErrInvalidInput = errors.New("invalid input")

func (uc *CreateDriverUseCase) Execute(ctx context.Context, input CreateDriverInput) error {
	if input.Name == "" || input.LastName == "" || input.Email == "" || input.Phone == "" || input.NumLicence == "" {
		return ErrInvalidInput
	}

	passwordRandom, err := randomString(10)
	if err != nil {
		return err
	}

	passwordHashed, err := utils.HashPassword(passwordRandom)
	if err != nil {
		return errors.New("error hashing password")
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

	user := entities.User{
		Name:     input.Name,
		LastName: input.LastName,
		Email:    input.Email,
		Password: passwordHashed,
		Role:     "driver",
	}

	if err := uc.userRepo.CreateUserTx(tx, &user); err != nil {
		return err
	}

	driver := entities.Driver{
		UserID:      user.ID,
		PhoneNumber: input.Phone,
		LicenseNo:   input.NumLicence,
	}
	if err := uc.driver.CreateDriverTx(tx, &driver); err != nil {
		return errors.New("error creating driver")
	}

	if err := uc.txProvider.CommitTx(ctx, tx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	committed = true

	return nil
}

func randomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	passwrod := make([]byte, length)
	for i := range passwrod {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		passwrod[i] = charset[num.Int64()]
	}

	return string(passwrod), nil

}
