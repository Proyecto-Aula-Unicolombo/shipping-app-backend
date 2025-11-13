package repository

import "shipping-app/internal/app/domain/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(id uint) (*entities.User, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]*entities.User, error)
	UpdateUser(user *entities.User) error
}
