package users

import (
	"errors"
	"shipping-app/internal/app/application/users"
	"shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UpdateUserRequest struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	NumLicence  string `json:"num_licence"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	NumLicence  string `json:"num_licence"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HandlerUser struct {
	createUserUseCase    *users.CreateUserUseCase
	getUserUseCase       *users.GetUser
	deleteUserUseCase    *users.DeleteUserUseCase
	listUsersUseCase     *users.ListUsers        // Tuyo (sin paginación)
	listUsersPaginatedUC *users.ListUsersUseCase // Del compañero (con paginación)
	updateUserUseCase    *users.UpdateUserUseCase
}

func NewHandlerUser(
	createUserUseCase *users.CreateUserUseCase,
	getUserUseCase *users.GetUser,
	deleteUserUseCase *users.DeleteUserUseCase,
	listUsersUseCase *users.ListUsers,
	listUsersPaginatedUC *users.ListUsersUseCase,
	updateUserUseCase *users.UpdateUserUseCase,
) *HandlerUser {
	return &HandlerUser{
		createUserUseCase:    createUserUseCase,
		getUserUseCase:       getUserUseCase,
		deleteUserUseCase:    deleteUserUseCase,
		listUsersUseCase:     listUsersUseCase,
		listUsersPaginatedUC: listUsersPaginatedUC,
		updateUserUseCase:    updateUserUseCase,
	}
}

func (h *HandlerUser) CreateUser(ctx fiber.Ctx) error {
	var req CreateUserRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid request body",
		})
	}
	input := users.CreateUserInput{
		Name:        req.Name,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
		Role:        req.Role,
		PhoneNumber: req.PhoneNumber,
		NumLicence:  req.NumLicence,
	}

	if err := h.createUserUseCase.Execute(ctx, input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully"})
}

func (h *HandlerUser) GetUser(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}

	user, err := h.getUserUseCase.Execute(uint(id))
	if err != nil {
		return h.handleGetUserError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario consultado exitosamente",
		"data":    user, // ← Ahora es UserOutput con info del driver
	})
}

func (h *HandlerUser) UpdateUser(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}

	var req UpdateUserRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Cuerpo de petición inválido",
		})
	}

	input := users.UpdateUserInput{
		ID:       uint(id),
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Driver: drivers.DriverUpdateDTO{
			PhoneNumber: req.PhoneNumber,
			LicenseNo:   req.NumLicence,
		},
	}

	err = h.updateUserUseCase.Execute(ctx, input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario actualizado correctamente",
	})
}

func (h *HandlerUser) DeleteUser(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "El ID debe ser un número válido",
		})
	}

	err = h.deleteUserUseCase.Execute(uint(id))
	if err != nil {
		return h.handleDeleteUserError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario eliminado correctamente",
	})
}

func (h *HandlerUser) ListUsersSimple(ctx fiber.Ctx) error {
	users, err := h.listUsersUseCase.Execute()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al listar usuarios",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  users,
		"total": len(users),
	})
}

func (h *HandlerUser) ListUsersPaginated(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)
	nameOrLastname := ctx.Query("name_or_last_name")
	role := ctx.Query("role")

	input := users.ListUserInput{
		Limit:          params.Limit,
		Offset:         params.Offset,
		NameOrLastname: nameOrLastname,
		Role:           role,
	}

	usersOuputs, total, err := h.listUsersPaginatedUC.Execute(input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al listar usuarios",
		})
	}
	if usersOuputs == nil {
		usersOuputs = []*users.ListUserOutput{}
	}

	response := utils.NewPaginationResponse(usersOuputs, int(total), params.Page, params.Limit)

	return ctx.JSON(response)
}

func (h *HandlerUser) handleGetUserError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, users.ErrUserNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "user_not_found",
			Message: "Usuario no registrado",
		})
	case errors.Is(err, users.ErrInvalidID):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_id",
			Message: "ID inválido",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al consultar usuario",
		})
	}
}

func (h *HandlerUser) handleDeleteUserError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, users.ErrUserNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "user_not_found",
			Message: "Usuario no registrado",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar usuario",
		})
	}
}

func (h *HandlerUser) handleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, users.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
	case errors.Is(err, users.ErrPasswordTooShort):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "password_too_short",
			Message: err.Error(),
		})
	case errors.Is(err, users.ErrInvalidRole):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_role",
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
