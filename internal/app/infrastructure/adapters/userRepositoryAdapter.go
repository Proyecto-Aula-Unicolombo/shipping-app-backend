package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) CreateUserTx(tx *sql.Tx, user *entities.User) error {
	query := `
		INSERT INTO users (name, lastname, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var err error
	if tx != nil {
		err = tx.QueryRow(query, user.Name, user.LastName, user.Email, user.Password, user.Role).Scan(&user.ID)
	}
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *UserRepositoryPostgres) GetUserByID(id uint) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, name, lastname, email, role FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) ListUsers(limit, offset int, NameOrLastname string) ([]*entities.User, error) {
	query := `SELECT id, name, lastname, email, role FROM users WHERE name ILIKE $1 OR lastname ILIKE $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, "%"+NameOrLastname+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer rows.Close()
	var users []*entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Role); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}
