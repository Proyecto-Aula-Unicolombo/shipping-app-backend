package handlerpackages

import (
	"errors"
	"log"
	usepackages "shipping-app/internal/app/application/UsePackages"
	"shipping-app/internal/app/domain/ports/repository"
	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

func (h *PackageHandler) handleErrorCancel(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usepackages.ErrToGetPackage):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "package_not_found",
			"message": "Package not found",
		})
	case errors.Is(err, usepackages.ErrToGetStatus):
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "status_retrieval_error",
			"message": "Could not retrieve package status",
		})
	case errors.Is(err, usepackages.ErrCannotCancel):
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "package_cannot_be_cancelled",
			"message": "The package cannot be cancelled in its current status",
		})

	default:
		// unknown/internal error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_server_error",
			"message": "An unexpected error occurred",
		})
	}
}

func (h *PackageHandler) handleErrorCreate(ctx fiber.Ctx, err error) error {
	log.Printf("Handling error: %v (type: %T)", err, err)

	// Verificar si es PackageConflictError (duplicate numpackage)
	var conflictErr *adapters.PackageConflictError
	if errors.As(err, &conflictErr) {
		response := fiber.Map{
			"error":      "package_already_exists",
			"message":    "A package with this number already exists",
			"numpackage": conflictErr.NumPackage,
		}

		if conflictErr.ExistingID > 0 {
			response["existing_id"] = conflictErr.ExistingID
		}

		return ctx.Status(fiber.StatusConflict).JSON(response)
	}

	// verificar si falta algun dato en la request
	var validationErr *usepackages.ValidationError
	if errors.As(err, &validationErr) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_input",
			"message": validationErr.Message,
			"fields":  validationErr.Fields,
		})
	}

	// otros casos de error
	switch {
	case errors.Is(err, usepackages.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_input",
			"message": err.Error(),
		})
	case errors.Is(err, usepackages.ErrRelatedEntityMustProvideID):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "related_entity_id_required",
			"message": err.Error(),
		})
	case errors.Is(err, usepackages.ErrRelatedEntityNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "related_entity_not_found",
			"message": err.Error(),
		})
	case errors.Is(err, usepackages.ErrBusinessRuleViolation):
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "business_rule_violation",
			"message": err.Error(),
		})
	default:
		// unknown/internal error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_server_error",
			"message": "could not create package",
		})
	}
}

func (h *PackageHandler) handleErrorConsultOrList(c fiber.Ctx, err error) error {
	switch err {
	case usepackages.ErrInvalidSearchCriteria:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_search_criteria",
			"message": err.Error(),
		})
	case usepackages.ErrAccessDenied:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "access_denied",
			"message": "You don't have permission to access this package",
		})
	case repository.ErrPackageNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "package_not_found",
			"message": "Package not found",
		})
	case usepackages.ErrGetRelatedEntities:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "related_entities_retrieval_error",
			"message": "Could not retrieve related entities for the package",
		})

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_server_error",
			"message": "An unexpected error occurred",
		})
	}
}
