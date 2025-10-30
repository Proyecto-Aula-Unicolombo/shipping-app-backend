package users

import (
	"errors"
	"shipping-app/internal/app/application/users"
	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HandlerUser struct {
	createUserUseCase *users.CreateUserUseCase
}

func NewHandlerUser(createUserUseCase *users.CreateUserUseCase) *HandlerUser {
	return &HandlerUser{createUserUseCase: createUserUseCase}
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
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.createUserUseCase.Execute(input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully"})
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
