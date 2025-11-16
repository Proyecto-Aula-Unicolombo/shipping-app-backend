package drivershandler

import (
	"errors"
	"shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

type CreateDriverRequest struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	NumLicence  string `json:"num_licence"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HandlerDrivers struct {
	createDriverUseCase *drivers.CreateDriverUseCase
}

func NewHandlerDrivers(createDriverUseCase *drivers.CreateDriverUseCase) *HandlerDrivers {
	return &HandlerDrivers{
		createDriverUseCase: createDriverUseCase,
	}
}

func (h *HandlerDrivers) CreateDriver(ctx fiber.Ctx) error {
	var req CreateDriverRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
		})
	}

	input := drivers.CreateDriverInput{
		Name:       req.Name,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.PhoneNumber,
		NumLicence: req.NumLicence,
	}

	if err := h.createDriverUseCase.Execute(ctx.Context(), input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Driver created successfully",
	})
}

func (h *HandlerDrivers) handleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, drivers.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})

	case errors.Is(err, adapters.ErrUserAlreadyExists):
		return ctx.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error:   "user_already_exists",
			Message: "A user with this email already exists",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Could not create user",
		})
	}
}
