package users

import (
	"context"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrUserNotFoundDelete = errors.New("Usuario no encontrado")
)

type DeleteUserUseCase struct {
	repo       repository.UserRepository
	driverRepo repository.DriverRepository
	txProvider repository.TxProvider
}

func NewDeleteUserUseCase(repo repository.UserRepository, driverRepo repository.DriverRepository, txProvider repository.TxProvider) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo, driverRepo: driverRepo, txProvider: txProvider}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return ErrUserNotFound
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

	if user.Role == "driver" {
		err = uc.driverRepo.DeleteDriverByUserIDTx(tx, id)
		if err != nil {
			return fmt.Errorf("delete driver: %w", err)
		}
	}
	err = uc.repo.DeleteUser(tx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	err = uc.txProvider.CommitTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	committed = true
	return nil
}
