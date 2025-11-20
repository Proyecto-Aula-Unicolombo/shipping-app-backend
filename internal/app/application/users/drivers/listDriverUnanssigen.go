package drivers

import "shipping-app/internal/app/domain/ports/repository"

type ListDriverUnassignedOutput struct {
	ID       uint
	Name     string
	LastName string
	License  string
}

type ListDriverUnassignedUseCase struct {
	driverRepo repository.DriverRepository
}

func NewListDriverUnassignedUseCase(driverRepo repository.DriverRepository) *ListDriverUnassignedUseCase {
	return &ListDriverUnassignedUseCase{driverRepo: driverRepo}
}

func (uc *ListDriverUnassignedUseCase) Execute() ([]*ListDriverUnassignedOutput, error) {
	drivers, err := uc.driverRepo.ListDriversUnassigned()
	if err != nil {
		return nil, err
	}

	driverOutput := []*ListDriverUnassignedOutput{}
	for _, driver := range drivers {
		driverOutput = append(driverOutput, &ListDriverUnassignedOutput{
			ID:       driver.ID,
			Name:     driver.User.Name,
			LastName: driver.User.LastName,
			License:  driver.LicenseNo,
		})
	}

	return driverOutput, nil
}
