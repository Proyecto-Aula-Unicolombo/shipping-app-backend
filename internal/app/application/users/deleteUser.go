package users

import (
	"errors"
	"shipping-app/internal/app/domain/ports/repository"
)

type DeleteUser struct {
	repo repository.UserRepository 
}

func NewDeleteUserUseCase(repo repository.UserRepository) *DeleteUser {
	return &DeleteUser{repo: repo} 
}

func (uc *DeleteUser) Execute(id uint) error { 
	
	_, err := uc.repo.GetUserByID(id) 
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	
	return uc.repo.DeleteUser(id)
}
