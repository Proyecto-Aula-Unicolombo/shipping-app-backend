package handlerpackages

import (
	usepackages "shipping-app/internal/app/application/UsePackages"
	related "shipping-app/internal/app/application/UsePackages/related"
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
	consultUC  *usepackages.ConsultPackageUseCase
}

func NewPackageHandler(createUC *usepackages.CreatePackageUseCase, cancelleUC *usepackages.CancelPackageUseCase, consultUC *usepackages.ConsultPackageUseCase) *PackageHandler {
	return &PackageHandler{createUC: createUC, cancelleUC: cancelleUC, consultUC: consultUC}
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

func (h *PackageHandler) ConsultPackageByNumPackage(ctx fiber.Ctx) error {
	numPackageStr := ctx.Params("numPackage")
	numPackage, err := strconv.ParseInt(numPackageStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_numpackage",
			"message": "NumPackage must be a valid number",
		})
	}

	// Obtener sender_id del contexto (puesto por el middleware APIKeyAuth)
	senderID, ok := ctx.Locals("sender_id").(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "unauthorized",
			"message": "Invalid sender authentication",
		})
	}

	input := usepackages.InputCheckAccess{
		Ctx:        ctx.Context(),
		NumPackage: &numPackage,
		AuthType:   "api_key",
		SenderID:   &senderID,
	}

	response, err := h.consultUC.Execute(input)
	if err != nil {
		return h.handleErrorCancel(ctx, err)
	}

	return ctx.JSON(response)
}

// Para UI (con JWT)
func (h *PackageHandler) ConsultPackageByID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_id",
			"message": "ID must be a valid number",
		})
	}

	packageID := uint(id)

	// Obtener información del usuario del contexto (puesto por JWTAuth)
	userRole, _ := ctx.Locals("user_role").(string)
	driverID, _ := ctx.Locals("driver_id").(*uint)

	input := usepackages.InputCheckAccess{
		Ctx:       ctx.Context(),
		PackageID: &packageID,
		AuthType:  "jwt",
		UserRole:  userRole,
		DriverID:  driverID,
	}

	response, err := h.consultUC.Execute(input)
	if err != nil {
		return h.handleErrorConsult(ctx, err)
	}

	return ctx.JSON(response)
}
