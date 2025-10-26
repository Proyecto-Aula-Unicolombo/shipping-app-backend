package users

import (
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/utils"
)

type CreateUserInput struct {
	Name     string
	LastName string
	Email    string
	Password string
	Role     string
}

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidRole      = errors.New("invalid role")
)

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (us *CreateUserUseCase) Execute(input CreateUserInput) error {
	if err := validateInput(input); err != nil {
		return err
	}
	passwordHashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return errors.New("error hashing password")
	}

	user := entities.User{
		Name:     input.Name,
		LastName: input.LastName,
		Email:    input.Email,
		Password: passwordHashed,
		Role:     input.Role,
	}

	return us.userRepo.CreateUser(&user)
}

func validateInput(input CreateUserInput) error {
	if input.Name == "" || input.LastName == "" || input.Email == "" || input.Password == "" || input.Role == "" {
		return ErrInvalidInput
	}

	if len(input.Password) < 8 {
		return ErrPasswordTooShort
	}

	validRoles := map[string]bool{
		"coord": true,
		"admin": true,
	}

	if !validRoles[input.Role] {
		return ErrInvalidRole
	}

	return nil
}
