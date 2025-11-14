package users

import (
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

var (
	ErrUserNotFound = errors.New("Usuario no registrado")
	ErrInvalidID    = errors.New("ID inválido")
)

type GetUser struct {
	repo repository.UserRepository
}

func NewGetUser(repo repository.UserRepository) *GetUser {
	return &GetUser{repo: repo}
}

func (uc *GetUser) Execute(id uint) (*entities.User, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Password = ""
	return user, nil
}