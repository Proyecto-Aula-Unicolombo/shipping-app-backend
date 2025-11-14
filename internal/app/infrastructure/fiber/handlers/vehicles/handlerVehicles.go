package vehicles

import (
	"errors"
	"shipping-app/internal/app/application/Vehicles"
	"shipping-app/internal/app/infrastructure/adapters"
    "strconv" 
	"github.com/gofiber/fiber/v3"
)

type CreateVehicleRequest struct {
	Plate       string `json:"plate"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Color       string `json:"color"`
	VehicleType string `json:"vehicleType"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HandlerVehicle struct {
	createVehicleUseCase *vehicles.CreateVehicleUseCase    
	getVehicleUseCase    *vehicles.GetVehicle
}

func NewHandlerVehicle(
	createVehicleUseCase *vehicles.CreateVehicleUseCase,
	getVehicleUseCase *vehicles.GetVehicle,  // ← AGREGAR
) *HandlerVehicle {
	return &HandlerVehicle{
		createVehicleUseCase: createVehicleUseCase,
		getVehicleUseCase:    getVehicleUseCase,  // ← AGREGAR
	}
}

func (h *HandlerVehicle) CreateVehicle(ctx fiber.Ctx) error {
	var req CreateVehicleRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid request body",
		})
	}

	input := vehicles.CreateVehicleInput{
		Plate:       req.Plate,
		Brand:       req.Brand,
		Model:       req.Model,
		Color:       req.Color,
		VehicleType: req.VehicleType,
	}

	if err := h.createVehicleUseCase.Execute(ctx, input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "vehicle created successfully",
	})
}

func (h *HandlerVehicle) GetVehicle(ctx fiber.Ctx) error {
	// 1. Extraer ID de la URL
	idParam := ctx.Params("id")
	
	// 2. Convertir a número
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}

	// 3. Llamar al caso de uso
	vehicle, err := h.getVehicleUseCase.Execute(uint(id))
	if err != nil {
		return h.handleGetVehicleError(ctx, err)
	}

	// 4. Devolver vehículo
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vehículo consultado exitosamente",
		"data":    vehicle,
	})
}

func (h *HandlerVehicle) handleGetVehicleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, vehicles.ErrVehicleNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "vehicle_not_found",
			Message: "Vehículo no registrado",
		})
	case errors.Is(err, vehicles.ErrInvalidID):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al consultar vehículo",
		})
	}
}

func (h *HandlerVehicle) handleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, vehicles.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})

	case errors.Is(err, adapters.ErrVehicleAlreadyExists):
		return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:   "vehicle_already_exists",
			Message: "A vehicle with this plate already exists",
		})

	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Could not create vehicle",
		})
	}
}
