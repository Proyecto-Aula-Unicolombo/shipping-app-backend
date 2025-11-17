package auth

import (
	"errors"
	authApp "shipping-app/internal/app/application/auth"
	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type AuthHandler struct {
	loginUseCase *authApp.LoginUseCase
}

func NewAuthHandler(loginUseCase *authApp.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *AuthHandler) Login(ctx fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	input := authApp.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUseCase.Execute(input)
	if err != nil {
		return h.handleLoginError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(output)
}

func (h *AuthHandler) handleLoginError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, authApp.ErrInvalidCredentials):
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Email or password is incorrect",
		})
	case errors.Is(err, authApp.ErrEmptyCredentials):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "empty_credentials",
			Message: "Email and password are required",
		})
	case errors.Is(err, adapters.ErrUserNotFound):
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Email or password is incorrect",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_server_error",
			Message: "Could not process login",
		})
	}
}
