package vehicles

import (
	"shipping-app/internal/app/domain/entities"
	"shipping-app/internal/app/domain/ports/repository"
)

type ListVehicles struct {
	repo repository.VehicleRepository
}

func NewListVehicles(repo repository.VehicleRepository) *ListVehicles {
	return &ListVehicles{repo: repo}
}

func (uc *ListVehicles) Execute() ([]*entities.Vehicle, error) {
	vehicles, err := uc.repo.GetAllVehicles()
	if err != nil {
		return nil, err
	}
	
	return vehicles, nil
}