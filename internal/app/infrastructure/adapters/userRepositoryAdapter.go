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

// NewUserRepositoryPostgres crea una instancia del repositorio
func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

// CreateUser inserta un nuevo usuario en la base de datos
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

// GetUserByID obtiene un usuario por su id
func (r *UserRepositoryPostgres) GetUserByID(id uint) (*entities.User, error) {
	var user entities.User

	query := `SELECT id, name, lastname, email, role FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Email,
		&user.Role,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("usuario no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryPostgres) DeleteUser(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil

	
}

func (r *UserRepositoryPostgres) GetAllUsers() ([]*entities.User, error) {
	query := `SELECT id, name, lastname, email, role FROM users ORDER BY id`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}
	defer rows.Close()
	
	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.LastName,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, &user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	
	return users, nil
}
func (r *UserRepositoryPostgres) UpdateUser(user *entities.User) error {
	query := `
		UPDATE users 
		SET name = $1, lastname = $2, email = $3, role = $4
		WHERE id = $5
	`

	res, err := r.db.Exec(
		query,
		user.Name,
		user.LastName,
		user.Email,
		user.Role,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}
