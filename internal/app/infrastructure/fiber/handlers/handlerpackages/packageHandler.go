package handlerpackages

import (
	"errors"
	usepackages "shipping-app/internal/app/application/UsePackages"
	"shipping-app/internal/app/application/UsePackages/related"

	"github.com/gofiber/fiber/v3"
)

type CreatePackageRequest struct {
	NumPackage           int64                              `json:"numpackage"`
	StartStatus          string                             `json:"startstatus"`
	DescriptionContent   *string                            `json:"descriptioncontent"`
	Weight               *float64                           `json:"weight"`
	Dimension            *float64                           `json:"dimension"`
	DeclaredValue        *float64                           `json:"declaredvalue"`
	TypePackage          *string                            `json:"typepackage"`
	IsFragile            bool                               `json:"is_fragile"`
	AddressPackage       *related.AdressPackageInput        `json:"addresspackage"`
	StatusDelivery       *related.StatusDeliveryInput       `json:"statusdelivery"`
	ComercialInformation *related.ComercialInformationInput `json:"comercialinformation"`
	Sender               *related.SenderInput               `json:"sender"`
	Receiver             *related.ReceiverInput             `json:"receiver"`
}

type PackageHandler struct {
	createUC *usepackages.CreatePackageUseCase
}

func NewPackageHandler(createUC *usepackages.CreatePackageUseCase) *PackageHandler {
	return &PackageHandler{createUC: createUC}
}

func (h *PackageHandler) CreatePackage(ctx fiber.Ctx) error {
	var req CreatePackageRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	input := usepackages.CreatePackageInput{
		NumPackage:           req.NumPackage,
		StartStatus:          req.StartStatus,
		DescriptionContent:   req.DescriptionContent,
		Weight:               req.Weight,
		Dimension:            req.Dimension,
		DeclaredValue:        req.DeclaredValue,
		TypePackage:          req.TypePackage,
		IsFragile:            req.IsFragile,
		AddressPackage:       req.AddressPackage,
		StatusDelivery:       req.StatusDelivery,
		ComercialInformation: req.ComercialInformation,
		Sender:               req.Sender,
		Receiver:             req.Receiver,
	}

	output, err := h.createUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         output.ID,
		"numpackage": output.NumPackage,
	})
}

func (h *PackageHandler) handleError(ctx fiber.Ctx, err error) error {
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
