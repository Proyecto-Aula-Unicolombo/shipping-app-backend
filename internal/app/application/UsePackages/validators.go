package usepackages

import (
	"errors"
	services "shipping-app/internal/app/domain/services/package"
)

var (
	ErrInvalidInput               = errors.New("invalid input")
	ErrRelatedEntityMustProvideID = errors.New("related entities must provide non-zero ID")
	ErrRelatedEntityNotFound      = errors.New("related entity not found")
	ErrBusinessRuleViolation      = errors.New("business rule violation")
)

type ValidationError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{
		Message: "invalid input",
		Fields:  fields,
	}
}

func ValidateCreateInput(input CreatePackageInput) error {
	fields := make(map[string]string)

	if input.NumPackage == "" {
		fields["numpackage"] = "numpackage is required"
	}

	if input.AddressPackage == nil {
		fields["addresspackage"] = "addresspackage is required"
	} else {
		if input.AddressPackage.Origin == "" {
			fields["addresspackage.origin"] = "origin is required"
		}
		if input.AddressPackage.Destination == "" {
			fields["addresspackage.destination"] = "destination is required"
		}
	}

	if input.ComercialInformation == nil {
		fields["comercialinformation"] = "comercialinformation is required"
	} else {
		if input.ComercialInformation.CostSending <= 0 {
			fields["comercialinformation.cost_sending"] = "cost_sending must be greater than 0"
		}
	}


	if input.Receiver == nil {
		fields["receiver"] = "receiver is required"
	} else {
		if input.Receiver.Name == "" {
			fields["receiver.name"] = "name is required"
		}
		if input.Receiver.LastName == "" {
			fields["receiver.last_name"] = "last_name is required"
		}
		if input.Receiver.PhoneNumber == "" {
			fields["receiver.phone_number"] = "phone_number is required"
		}
		if input.Receiver.Email == "" {
			fields["receiver.email"] = "email is required"
		}
	}

	if len(fields) > 0 {
		return NewValidationError(fields)
	}

	return nil
}

func ValidateBusinessRules(domainSvc *services.ValidatePackageService, input CreatePackageInput) error {
	if domainSvc == nil {
		return nil
	}
	// Dimension ahora es string descriptivo, no se valida numéricamente
	if err := domainSvc.ValidatePackageBusinessRules(input.Weight, nil, input.DeclaredValue, input.IsFragile); err != nil {
		return err
	}
	return nil
}
