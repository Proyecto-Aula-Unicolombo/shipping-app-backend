package vehicles

import (
	"errors"
	"shipping-app/internal/app/application/Vehicles"
	"shipping-app/internal/app/infrastructure/adapters"

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
}

func NewHandlerVehicle(createVehicleUseCase *vehicles.CreateVehicleUseCase) *HandlerVehicle {
	return &HandlerVehicle{createVehicleUseCase: createVehicleUseCase}
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
