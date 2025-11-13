package users

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

type UpdateUserUseCase struct {
	repo repository.UserRepository
}

func NewUpdateUserUseCase(repo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

type UpdateUserInput struct {
	ID       uint
	Name     string
	LastName string
	Email    string
	Role     string
}

func (uc *UpdateUserUseCase) Execute(input UpdateUserInput) error {
	if input.ID == 0 {
		return errors.New("ID inválido")
	}

	existingUser, err := uc.repo.GetUserByID(input.ID)
	if err != nil {
		return errors.New("usuario no encontrado")
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
	if input.Role != "" {
		existingUser.Role = input.Role
	}

	
	return uc.repo.UpdateUser(existingUser)
}
