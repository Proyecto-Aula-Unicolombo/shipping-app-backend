package tracking

import (
	"shipping-app/internal/app/application/tracking"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type TrackPackageByNumRequest struct {
	NumPackage string `json:"num_package"`
	ReceiverID *uint  `json:"receiver_id"`
}

type RegisterStopRequest struct {
	OrderID     uint    `json:"order_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TypeStop    string  `json:"type_stop"`
	Description *string `json:"description"`
	Evidence    *string `json:"evidence"`
}

type TrackingHandler struct {
	trackPackageUC *tracking.TrackPackageUseCase
	registerStopUC *tracking.RegisterStopUseCase
}

func NewTrackingHandler(
	trackPackageUC *tracking.TrackPackageUseCase,
	registerStopUC *tracking.RegisterStopUseCase,
) *TrackingHandler {
	return &TrackingHandler{
		trackPackageUC: trackPackageUC,
		registerStopUC: registerStopUC,
	}
}

// TrackPackageByNumber - Para uso público del destinatario
func (h *TrackingHandler) TrackPackageByNumber(ctx fiber.Ctx) error {
	numPackage := ctx.Query("num_package")
	if numPackage == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "missing_parameter",
			"message": "num_package query parameter is required",
		})
	}

	// ReceiverID opcional desde query o token (si está autenticado)
	var receiverID *uint
	receiverIDStr := ctx.Query("receiver_id")
	if receiverIDStr != "" {
		id, err := strconv.ParseUint(receiverIDStr, 10, 32)
		if err == nil {
			receiverIDUint := uint(id)
			receiverID = &receiverIDUint
		}
	}

	input := tracking.TrackPackageInput{
		NumPackage: &numPackage,
		ReceiverID: receiverID,
	}

	response, err := h.trackPackageUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data": response,
	})
}

// TrackPackageByID - Para uso interno
func (h *TrackingHandler) TrackPackageByID(ctx fiber.Ctx) error {
	packageIDStr := ctx.Params("packageId")
	packageID, err := strconv.ParseUint(packageIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_package_id",
			"message": "Package ID must be a valid number",
		})
	}

	packageIDUint := uint(packageID)
	input := tracking.TrackPackageInput{
		PackageID: &packageIDUint,
	}

	response, err := h.trackPackageUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data": response,
	})
}

func (h *TrackingHandler) RegisterStop(ctx fiber.Ctx) error {
	var req RegisterStopRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	input := tracking.RegisterStopInput{
		OrderID:     req.OrderID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		TypeStop:    req.TypeStop,
		Description: req.Description,
		Evidence:    req.Evidence,
	}

	output, err := h.registerStopUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "stop registered successfully",
		"data":    output,
	})
}

func (h *TrackingHandler) handleError(ctx fiber.Ctx, err error) error {
	switch err {
	case tracking.ErrInvalidTrackingInput, tracking.ErrInvalidStopInput, tracking.ErrInvalidStopType:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_input",
			"message": err.Error(),
		})
	case tracking.ErrUnauthorizedAccess:
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "forbidden",
			"message": err.Error(),
		})
	case tracking.ErrOrderNotFound:
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": err.Error(),
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "An error occurred while processing the request",
		})
	}
}
