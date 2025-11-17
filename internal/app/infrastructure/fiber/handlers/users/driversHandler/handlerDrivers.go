package drivershandler

import (
	"errors"
	"shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/utils"
	"strconv"

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
	createDriverUseCase       *drivers.CreateDriverUseCase
	ListDriversUseCase        *drivers.ListDriverUseCase
	GetDriverByIdUseCase      *drivers.GetDriversByIdUseCase
	UpdateStatusDriverUseCase *drivers.UpdateStatusDriverUseCase
}

func NewHandlerDrivers(createDriverUseCase *drivers.CreateDriverUseCase, listDriversUseCase *drivers.ListDriverUseCase, getDriverByIdUseCase *drivers.GetDriversByIdUseCase, updateStatusDriverUseCase *drivers.UpdateStatusDriverUseCase) *HandlerDrivers {
	return &HandlerDrivers{
		createDriverUseCase:       createDriverUseCase,
		ListDriversUseCase:        listDriversUseCase,
		GetDriverByIdUseCase:      getDriverByIdUseCase,
		UpdateStatusDriverUseCase: updateStatusDriverUseCase,
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

func (h *HandlerDrivers) ListDrivers(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)
	nameORLastName := ctx.Query("name_or_lastname")

	input := drivers.ListDriverInput{
		Limit:          params.Limit,
		Offset:         params.Offset,
		NameOrLastName: nameORLastName,
	}

	driversOutput, total, err := h.ListDriversUseCase.Execute(input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al listar drivers",
		})
	}

	if driversOutput == nil {
		driversOutput = []*drivers.ListDriverOutput{}
	}

	response := utils.NewPaginationResponse(driversOutput, int(total), params.Page, params.Limit)

	return ctx.JSON(response)
}

func (h *HandlerDrivers) GetDriverByID(ctx fiber.Ctx) error {
	idstr := ctx.Params("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}
	driverOutput, err := h.GetDriverByIdUseCase.Execute(ctx.Context(), uint(id))
	if err != nil {
		if errors.Is(err, drivers.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "not_found",
				Message: "Driver not found",
			})
		} else {
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "internal_server_error",
				Message: "Could not retrieve driver",
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(driverOutput)
}

func (h *HandlerDrivers) UpdateStatusDriver(ctx fiber.Ctx) error {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	idstr := ctx.Params("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}

	if err := h.UpdateStatusDriverUseCase.Execute(uint(id), req.IsActive); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Driver status updated successfully",
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
