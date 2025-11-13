package users

import (
	"errors"
	"shipping-app/internal/app/application/users"
	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

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

	case errors.Is(err, users.ErrNoUsersFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "users_no_found",
			Message: "Users not found",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Could not create user",
		})
	}
}
