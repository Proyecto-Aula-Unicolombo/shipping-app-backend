package drivers

import "shipping-app/internal/app/domain/ports/repository"

type ListDriverInput struct {
	Limit  int
	Offset int

	NameOrLastName string
}

type ListDriverOutput struct {
	ID       uint
	Name     string
	LastName string
	NumOrder uint
	IsActive bool
}

type ListDriverUseCase struct {
	driverRepo repository.DriverRepository
	userRepo   repository.UserRepository
}

func NewListDriverUseCase(driverRepo repository.DriverRepository, userRepo repository.UserRepository) *ListDriverUseCase {
	return &ListDriverUseCase{
		driverRepo: driverRepo,
		userRepo:   userRepo,
	}
}

func (uc *ListDriverUseCase) Execute(input ListDriverInput) ([]*ListDriverOutput, int64, error) {
	total, err := uc.driverRepo.CountDrivers(input.NameOrLastName)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*ListDriverOutput{}, 0, nil
	}

	drivers, err := uc.driverRepo.ListDrivers(input.Limit, input.Offset, input.NameOrLastName)
	if err != nil {
		return nil, 0, err
	}

	var outputs []*ListDriverOutput

	for _, driver := range drivers {
		outputs = append(outputs, &ListDriverOutput{
			ID:       driver.ID,
			Name:     driver.User.Name,
			LastName: driver.User.LastName,
			NumOrder: driver.NumOrder,
			IsActive: driver.IsActive,
		})
	}

	return outputs, total, nil

}
