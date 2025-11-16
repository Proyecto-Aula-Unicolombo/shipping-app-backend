package users

import (
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListUserInput struct {
	Limit  int
	Offset int

	NameOrLastname string
	Role           string
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

func (uc *ListUsersUseCase) Execute(input ListUserInput) ([]*ListUserOutput, int64, error) {
	total, err := uc.userRepo.CountUsers(input.NameOrLastname, input.Role)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*ListUserOutput{}, 0, nil
	}
	var users []*entities.User

	users, err = uc.userRepo.ListUsers(input.Limit, input.Offset, input.NameOrLastname, input.Role)
	if err != nil {
		return nil, 0, err
	}

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
