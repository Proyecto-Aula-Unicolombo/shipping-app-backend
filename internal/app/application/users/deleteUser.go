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
	// 1. Validar ID
	if id == 0 {
		return ErrInvalidID
	}

	// 2. Verificar que el usuario existe
	_, err := uc.repo.GetUserByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	// 3. Eliminar usuario
	return uc.repo.DeleteUser(id)
}