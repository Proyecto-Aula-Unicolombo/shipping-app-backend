package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrUserNotFound = errors.New("user not found")
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
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryPostgres) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, name, lastname, email, password, role FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	return &user, nil
}

// ListUsers del compañero (con paginación y búsqueda)
func (r *UserRepositoryPostgres) ListUsers(limit, offset int, nameOrLastname, role string) ([]*entities.User, error) {
	query := `SELECT id, name, lastname, email, role FROM users WHERE 1=1`
	args := []interface{}{}
	argPosition := 1

	if nameOrLastname != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR lastname ILIKE $%d)", argPosition, argPosition)
		args = append(args, "%"+nameOrLastname+"%")
		argPosition++
	}

	if role != "" && role != "all" {
		query += fmt.Sprintf(" AND role = $%d", argPosition)
		args = append(args, role)
		argPosition++
	}

	query += " ORDER BY id LIMIT $" + fmt.Sprint(argPosition) + " OFFSET $" + fmt.Sprint(argPosition+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
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

func (r *UserRepositoryPostgres) DeleteUser(tx *sql.Tx, id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := tx.Exec(query, id)
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

func (r *UserRepositoryPostgres) UpdateUser(tx *sql.Tx, user *entities.User) error {
	query := `
		UPDATE users 
		SET name = $1, lastname = $2, email = $3, role = $4
		WHERE id = $5
	`

	res, err := tx.Exec(
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

func (r *UserRepositoryPostgres) CountUsers(nameOrLastname, role string) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE (name ILIKE $1 OR lastname ILIKE $1) AND role ILIKE $2`
	var count int64
	err := r.db.QueryRow(query, "%"+nameOrLastname+"%", "%"+role+"%").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}
