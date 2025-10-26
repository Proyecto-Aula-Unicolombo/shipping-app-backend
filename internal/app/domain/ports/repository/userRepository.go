package repository

import "shipping-app/internal/app/domain/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(id uint) (*entities.User, error)
}
