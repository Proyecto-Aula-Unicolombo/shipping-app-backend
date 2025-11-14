package users

import (
	"errors"
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListUserInput struct {
	Limit  int
	Offset int

	NameOrLastname string
}

type ListUserOutput struct {
	ID       uint
	Name     string
	LastName string
	Email    string
	Role     string
}

type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

var ErrNoUsersFound = errors.New("no users found")

func (uc *ListUsersUseCase) Execute(input ListUserInput) ([]*ListUserOutput, int64, error) {
	var total int64
	var users []*entities.User
	var err error

	users, err = uc.userRepo.ListUsers(input.Limit, input.Offset, input.NameOrLastname)
	if err != nil {
		return nil, 0, err
	}
	if len(users) == 0 {
		return nil, 0, ErrNoUsersFound
	}
	total = int64(len(users))

	var outputs []*ListUserOutput
	for _, user := range users {
		outputs = append(outputs, &ListUserOutput{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Role:     user.Role,
		})
	}

	return outputs, total, nil
}
