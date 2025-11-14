package repository

import (
	"database/sql"
	"shipping-app/internal/app/domain/entities"
)

type UserRepository interface {
	CreateUserTx(tx *sql.Tx, user *entities.User) error
	GetUserByID(id uint) (*entities.User, error)
	DeleteUser(id uint) error
	UpdateUser(user *entities.User) error
	GetAllUsers() ([]*entities.User, error)
	ListUsers(limit, offset int, NameOrLastname string) ([]*entities.User, error)
}