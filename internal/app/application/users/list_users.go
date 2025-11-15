package users

import (
    "shipping-app/internal/app/domain/entities"
    "shipping-app/internal/app/domain/ports/repository"
)

type ListUsers struct {
    repo repository.UserRepository
}

func NewListUsers(repo repository.UserRepository) *ListUsers {
    return &ListUsers{repo: repo}
}

func (uc *ListUsers) Execute() ([]*entities.User, error) {
    users, err := uc.repo.GetAllUsers()
    if err != nil {
        return nil, err
    }
    
    for _, user := range users {
        user.Password = ""
    }
    
    return users, nil
}