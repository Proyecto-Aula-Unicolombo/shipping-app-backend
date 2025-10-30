package usepackages

import (
	"fmt"
	services "shipping-app/internal/app/domain/services/package"
)

func ValidateCreateInput(input CreatePackageInput) error {
	if input.NumPackage == 0 {
		return ErrInvalidInput
	}
	if input.StartStatus == "" {
		return ErrInvalidInput
	}
	if input.AddressPackage == nil || input.StatusDelivery == nil || input.ComercialInformation == nil ||
		input.Sender == nil || input.Receiver == nil {
		return ErrInvalidInput
	}
	return nil
}

func ValidateBusinessRules(domainSvc *services.ValidatePackageService, input CreatePackageInput) error {
	if domainSvc == nil {
		return nil
	}
	if err := domainSvc.ValidatePackageBusinessRules(input.Weight, input.Dimension, input.DeclaredValue, input.IsFragile); err != nil {
		return fmt.Errorf("%w: %v", ErrBusinessRuleViolation, err)
	}
	return nil
}
