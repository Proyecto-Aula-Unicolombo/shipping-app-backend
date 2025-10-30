package services

import "fmt"

type ValidatePackageService struct {
	// si necesitamos dependencias, las agregamos aquí ej: repositorios, config, etc.
}

func NewValidatePackageService() *ValidatePackageService {
	return &ValidatePackageService{}
}

var (
	ErrWeightExceeded    = fmt.Errorf("weight exceeds allowed limit")
	ErrInvalidDimensions = fmt.Errorf("invalid dimensions")
	ErrInvalidValue      = fmt.Errorf("invalid declared value")
	ErrInvalidPackage    = fmt.Errorf("invalid package data")
)

// ValidatePackageBusinessRules valida reglas de negocio básicas del paquete.
// - weight y dimension pueden ser nil (no provistos), pero si provistos deben ser > 0.
// - si isFragile == true, se aplica un límite máximo de peso (ej: 30.0).
// - declaredValue si provisto no debe ser negativo.
func (s *ValidatePackageService) ValidatePackageBusinessRules(weight *float64, dimension *float64, declaredValue *float64, isFragile bool) error {
	if weight != nil {
		if *weight <= 0 {
			return ErrInvalidDimensions
		}
		if isFragile && *weight > 30.0 {
			return ErrWeightExceeded
		}
	}
	if dimension != nil {
		if *dimension <= 0 {
			return ErrInvalidDimensions
		}
	}

	if declaredValue != nil {
		if *declaredValue < 0 {
			return ErrInvalidValue
		}
	}
	return nil
}
