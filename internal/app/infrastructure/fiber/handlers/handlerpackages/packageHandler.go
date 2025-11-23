package handlerpackages

import (
	"fmt"
	usepackages "shipping-app/internal/app/application/UsePackages"
	related "shipping-app/internal/app/application/UsePackages/related"
	"shipping-app/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CreatePackageRequest struct {
	NumPackage           string                             `json:"numpackage"`
	DescriptionContent   *string                            `json:"descriptioncontent"`
	Weight               *float64                           `json:"weight"`
	Dimension            *string                            `json:"dimension"`
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
	createUC               *usepackages.CreatePackageUseCase
	cancelleUC             *usepackages.CancelPackageUseCase
	consultUC              *usepackages.ConsultPackageUseCase
	listPkgUC              *usepackages.ListPackagesUseCase
	listPkgToCreateOrderUC *usepackages.ListPackagesToCreateOrderUseCase
}

func NewPackageHandler(
	createUC *usepackages.CreatePackageUseCase,
	cancelleUC *usepackages.CancelPackageUseCase,
	consultUC *usepackages.ConsultPackageUseCase,
	listPkgUC *usepackages.ListPackagesUseCase,
	listPkgToCreateOrderUC *usepackages.ListPackagesToCreateOrderUseCase,
) *PackageHandler {
	return &PackageHandler{
		createUC:               createUC,
		cancelleUC:             cancelleUC,
		consultUC:              consultUC,
		listPkgUC:              listPkgUC,
		listPkgToCreateOrderUC: listPkgToCreateOrderUC,
	}
}

func (h *PackageHandler) CreatePackage(ctx fiber.Ctx) error {
	var reqArray []CreatePackageRequest
	if err := ctx.Bind().Body(&reqArray); err == nil && len(reqArray) > 0 {
		return h.createMultiplePackages(ctx, reqArray)
	}

	// Si falla, intentar como objeto único
	var req CreatePackageRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request format",
		})
	}

	// Es un solo paquete
	return h.createSinglePackage(ctx, req)
}

func (h *PackageHandler) createSinglePackage(ctx fiber.Ctx, req CreatePackageRequest) error {
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

func (h *PackageHandler) createMultiplePackages(ctx fiber.Ctx, requests []CreatePackageRequest) error {
	inputs := make([]usepackages.CreatePackageInput, len(requests))
	for i, req := range requests {
		inputs[i] = usepackages.CreatePackageInput{
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
	}

	outputs, err := h.createUC.ExecuteBulk(ctx.Context(), inputs)
	if err != nil {
		return h.handleErrorCreate(ctx, err)
	}

	// Construir respuesta con todos los paquetes creados
	packagesCreated := make([]fiber.Map, len(outputs))
	for i, output := range outputs {
		packagesCreated[i] = fiber.Map{
			"id":         output.ID,
			"numpackage": output.NumPackage,
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  fmt.Sprintf("%d packages created successfully", len(outputs)),
		"packages": packagesCreated,
	})
}

func (h *PackageHandler) DeletePackage(ctx fiber.Ctx) error {
	numPackage := ctx.Params("numPackage")

	if numPackage == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request: numPackage is required",
		})
	}

	err := h.cancelleUC.Execute(ctx.Context(), numPackage)
	if err != nil {
		return h.handleErrorCancel(ctx, err)
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *PackageHandler) ConsultPackageByNumPackage(ctx fiber.Ctx) error {
	numPackage := ctx.Params("numPackage")
	if numPackage == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_numpackage",
			"message": "NumPackage is required",
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

	input := usepackages.CheckAccessInput{
		Ctx:        ctx.Context(),
		NumPackage: &numPackage,
		AuthType:   "api_key",
		SenderID:   &senderID,
	}

	response, err := h.consultUC.Execute(input)
	if err != nil {
		return h.handleErrorConsultOrList(ctx, err)
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

	input := usepackages.CheckAccessInput{
		Ctx:       ctx.Context(),
		PackageID: &packageID,
		AuthType:  "jwt",
		UserRole:  userRole,
	}

	response, err := h.consultUC.Execute(input)
	if err != nil {
		return h.handleErrorConsultOrList(ctx, err)
	}

	return ctx.JSON(response)
}

func (h *PackageHandler) ListPackages(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)

	authType := getAuthType(ctx)
	userRole, _ := ctx.Locals("user_role").(string)
	senderID, _ := ctx.Locals("sender_id").(uint)

	input := usepackages.ListPackagesInput{
		Ctx:      ctx.Context(),
		Limit:    params.Limit,
		Offset:   params.Offset,
		AuthType: authType,
		UserRole: userRole,
	}

	if senderID != 0 {
		input.SenderID = &senderID
	}
	packages, total, err := h.listPkgUC.Execute(input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error":   "internal_error",
			"Message": "Could not get packages",
		})
	}

	if packages == nil {
		packages = []*usepackages.ResponsePackage{}
	}

	response := utils.NewPaginationResponse(packages, int(total), params.Page, params.Limit)

	return ctx.JSON(response)
}

func (h *PackageHandler) ListPackagesToCreateOrder(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)

	input := usepackages.ListPackagesToCreateOrderInput{
		Ctx:    ctx.Context(),
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	packages, total, err := h.listPkgToCreateOrderUC.Execute(input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error":   "internal_error",
			"Message": "Could not get packages",
		})
	}

	if packages == nil {
		packages = []*usepackages.ListPackagesToCreateOrderOutput{}
	}

	response := utils.NewPaginationResponse(packages, int(total), params.Page, params.Limit)

	return ctx.JSON(response)
}

func getAuthType(c fiber.Ctx) string {
	// Verificar si hay JWT
	if c.Locals("user_role") != nil {
		return "jwt"
	}

	// Verificar si hay API Key
	if c.Locals("sender_id") != nil {
		return "api_key"
	}

	return ""
}
