package usepackages

import (
	"context"
	"shipping-app/internal/app/domain/entities"
)

type InputCheckAccess struct {
	Ctx context.Context

	PackageID  *uint  `json:"package_id,omitempty"`  // Para UI (interno)
	NumPackage *int64 `json:"num_package,omitempty"` // Para API Key (externo)

	AuthType string `json:"-"` // "jwt" o "api_key"
	UserRole string `json:"-"` // "coordinator", "driver"
	DriverID *uint  `json:"-"` // ID del conductor (para filtrar)
	SenderID *uint  `json:"-"` // ID del sender (para filtrar)

}

// checkAccess verifica si el usuario/sender tiene acceso al paquete
func CheckAccess(pkg *entities.Package, input InputCheckAccess) error {
	switch input.AuthType {
	case "api_key":
		if input.SenderID != nil && pkg.SenderID != *input.SenderID {
			return ErrAccessDenied
		}

	case "jwt":
		switch input.UserRole {
		case "coordinator":
			return nil

		case "driver":
			return nil
		default:
			return ErrAccessDenied
		}
	}

	return nil
}
