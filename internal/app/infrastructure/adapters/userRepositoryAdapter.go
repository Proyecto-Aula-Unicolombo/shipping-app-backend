package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

// DeleteUser implements repository.UserRepository.
func (r *UserRepositoryPostgres) DeleteUser(id uint) error {
	panic("unimplemented")
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) CreateUser(user *entities.User) error {
	query := `
		INSERT INTO users (name, lastname, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		user.Name,
		user.LastName,
		user.Email,
		user.Password,
		user.Role,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *UserRepositoryPostgres) GetUserByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT id, name, lastName, email, role FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
