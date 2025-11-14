package users

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrUserNotFoundDelete = errors.New("Usuario no encontrado")
)

type DeleteUserUseCase struct {
	repo repository.UserRepository
}

func NewDeleteUserUseCase(repo repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

func (uc *DeleteUserUseCase) Execute(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	_, err := uc.repo.GetUserByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return uc.repo.DeleteUser(id)
}