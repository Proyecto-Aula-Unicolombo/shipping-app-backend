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
	trackPackageUC     *tracking.TrackPackageUseCase
	registerStopUC     *tracking.RegisterStopUseCase
	listIncidentsUC    *tracking.ListIncidentsUseCase
	listActiveOrdersUC *tracking.ListActiveOrdersUseCase
}

func NewTrackingHandler(
	trackPackageUC *tracking.TrackPackageUseCase,
	registerStopUC *tracking.RegisterStopUseCase,
	listIncidentsUC *tracking.ListIncidentsUseCase,
	listActiveOrdersUC *tracking.ListActiveOrdersUseCase,
) *TrackingHandler {
	return &TrackingHandler{
		trackPackageUC:     trackPackageUC,
		registerStopUC:     registerStopUC,
		listIncidentsUC:    listIncidentsUC,
		listActiveOrdersUC: listActiveOrdersUC,
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

func (h *TrackingHandler) ListIncidents(ctx fiber.Ctx) error {
	// Parse query parameters for filters
	var status *string
	statusStr := ctx.Query("status")
	if statusStr != "" {
		status = &statusStr
	}

	var driverID *uint
	driverIDStr := ctx.Query("driver_id")
	if driverIDStr != "" {
		id, err := strconv.ParseUint(driverIDStr, 10, 32)
		if err == nil {
			driverIDUint := uint(id)
			driverID = &driverIDUint
		}
	}

	var orderID *uint
	orderIDStr := ctx.Query("order_id")
	if orderIDStr != "" {
		id, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err == nil {
			orderIDUint := uint(id)
			orderID = &orderIDUint
		}
	}

	// Pagination
	limit := 50 // default
	limitStr := ctx.Query("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	offsetStr := ctx.Query("offset")
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	input := tracking.ListIncidentsInput{
		Status:   status,
		DriverID: driverID,
		OrderID:  orderID,
		Limit:    limit,
		Offset:   offset,
	}

	incidents, err := h.listIncidentsUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data":  incidents,
		"count": len(incidents),
	})
}

func (h *TrackingHandler) ListActiveOrders(ctx fiber.Ctx) error {
	// Este endpoint devuelve todas las órdenes activas con su última ubicación
	// Solo accesible para admin/coordinador
	orders, err := h.listActiveOrdersUC.Execute(ctx.Context())
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data":  orders,
		"count": len(orders),
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
