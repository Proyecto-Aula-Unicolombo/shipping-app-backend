package users

import (
	"shipping-app/internal/app/application/users"
	"shipping-app/internal/utils"

	"github.com/gofiber/fiber/v3"
)

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
	createUserUseCase *users.CreateUserUseCase
	listUsersUseCase  *users.ListUsersUseCase
}

func NewHandlerUser(createUserUseCase *users.CreateUserUseCase, listUsersUseCase *users.ListUsersUseCase) *HandlerUser {
	return &HandlerUser{createUserUseCase: createUserUseCase, listUsersUseCase: listUsersUseCase}
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

func (h *HandlerUser) ListUsers(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)
	nameOrLastname := ctx.Query("name_or_last_name")

	input := users.ListUserInput{
		Limit:          params.Limit,
		Offset:         params.Offset,
		NameOrLastname: nameOrLastname,
	}
	users, total, err := h.listUsersUseCase.Execute(input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	response := utils.NewPaginationResponse(users, int(total), params.Page, params.Limit)

	return ctx.JSON(response)
}
