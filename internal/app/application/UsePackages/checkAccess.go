package usepackages

import (
	"context"
	"errors"
	"shipping-app/internal/app/domain/entities"
)

var (
	ErrInvalidAuthType        = errors.New("invalid authentication type")
	ErrMissingAuthCredentials = errors.New("missing authentication credentials")
)

type CheckAccessInput struct {
	Ctx        context.Context
	NumPackage *string
	PackageID  *uint
	AuthType   string `json:"-"` // "jwt" o "api_key"
	UserRole   string `json:"-"` // "coordinator", "driver"
	SenderID   *uint  `json:"-"` // ID del sender (para filtrar)

}

// checkAccess verifica si el usuario/sender tiene acceso al paquete
func CheckAccess(pkg *entities.Package, input CheckAccessInput) error {
	switch input.AuthType {
	case "api_key":
		return checkAPIKeyAccess(pkg, input)

	case "jwt":
		return nil

	default:
		return ErrInvalidAuthType
	}
}

func checkAPIKeyAccess(pkg *entities.Package, input CheckAccessInput) error {
	if input.SenderID == nil {
		return ErrMissingAuthCredentials
	}

	// Sender solo puede acceder a sus propios paquetes
	if pkg.SenderID != *input.SenderID {
		return ErrAccessDenied
	}

	return nil
}

func CheckBulkAccess(packages []*entities.Package, input CheckAccessInput) ([]*entities.Package, error) {
	var allowedPackages []*entities.Package

	for _, pkg := range packages {
		if err := CheckAccess(pkg, input); err == nil {
			allowedPackages = append(allowedPackages, pkg)
		}
	}

	return allowedPackages, nil
}
