package handlerpackages

import (
	usepackages "shipping-app/internal/app/application/UsePackages"
	"shipping-app/internal/app/application/UsePackages/related"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CreatePackageRequest struct {
	NumPackage           int64                              `json:"numpackage"`
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
	createUC   *usepackages.CreatePackageUseCase
	cancelleUC *usepackages.CancelPackageUseCase
}

func NewPackageHandler(createUC *usepackages.CreatePackageUseCase, cancelleUC *usepackages.CancelPackageUseCase) *PackageHandler {
	return &PackageHandler{createUC: createUC, cancelleUC: cancelleUC}
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
		return h.handleErrorCreate(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "package created successfully",
		"id":         output.ID,
		"numpackage": output.NumPackage,
	})
}

func (h *PackageHandler) DeletePackage(ctx fiber.Ctx) error {
	numPackageSTR := ctx.Params("numPackage")

	numPackage, err := strconv.Atoi(numPackageSTR)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	err = h.cancelleUC.Execute(ctx.Context(), int64(numPackage))
	if err != nil {
		return h.handleErrorCancel(ctx, err)
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
