package vehicles

import "shipping-app/internal/app/domain/ports/repository"

type ListUnassignedVehiclesOutput struct {
	ID          uint
	Plate       string
	Brand       string
	Model       string
	VehicleType string
}

type ListUnassignedVehiclesUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewListUnassignedVehiclesUseCase(vehicleRepo repository.VehicleRepository) *ListUnassignedVehiclesUseCase {
	return &ListUnassignedVehiclesUseCase{vehicleRepo: vehicleRepo}
}

func (uc *ListUnassignedVehiclesUseCase) Execute() ([]*ListUnassignedVehiclesOutput, error) {
	vehicles, err := uc.vehicleRepo.ListVehiclesUnassigned()
	if err != nil {
		return nil, err
	}
	var outputs []*ListUnassignedVehiclesOutput
	for _, v := range vehicles {
		outputs = append(outputs, &ListUnassignedVehiclesOutput{
			ID:          v.ID,
			Plate:       v.Plate,
			Brand:       v.Brand,
			Model:       v.Model,
			VehicleType: v.VehicleType,
		})
	}
	return outputs, nil
}
