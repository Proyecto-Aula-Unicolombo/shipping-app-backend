package vehicles

import (
	"shipping-app/internal/app/domain/ports/repository"
)

type ListVehiclesInput struct {
	Limit  int
	Offset int

	PlateBrandOrModel string
}

type ListVehiclesOutput struct {
	ID             uint
	Plate          string
	Brand          string
	Model          string
	Color          string
	VehicleType    string
	DriverName     string
	DriverLastName string
}

type ListVehicles struct {
	repo repository.VehicleRepository
}

func NewListVehicles(repo repository.VehicleRepository) *ListVehicles {
	return &ListVehicles{repo: repo}
}

func (uc *ListVehicles) Execute(input ListVehiclesInput) ([]*ListVehiclesOutput, int64, error) {
	count, err := uc.repo.CountVehicles(input.PlateBrandOrModel)
	if err != nil {
		return nil, 0, err
	}

	vehicles, err := uc.repo.ListVehicles(input.Limit, input.Offset, input.PlateBrandOrModel)
	if err != nil {
		return nil, 0, err
	}
	outputs := []*ListVehiclesOutput{}
	for _, v := range vehicles {
		outputs = append(outputs, &ListVehiclesOutput{
			ID:             v.ID,
			Plate:          v.Plate,
			Brand:          v.Brand,
			Model:          v.Model,
			Color:          v.Color,
			VehicleType:    v.VehicleType,
			DriverName:     v.AssignedDriverName,
			DriverLastName: v.AssignedDriverLastName,
		})
	}

	return outputs, count, nil
}
