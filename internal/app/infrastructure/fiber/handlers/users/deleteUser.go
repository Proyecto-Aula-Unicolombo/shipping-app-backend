package users

import (
	"strconv"
	"shipping-app/internal/app/application/users"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	deleteUserUseCase *users.DeleteUser  // ← Corregido nombre
}

func NewUserHandler(deleteUC *users.DeleteUser) *UserHandler {  // ← Corregido
	return &UserHandler{deleteUserUseCase: deleteUC}
}

func (h *UserHandler) DeleteUser(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)  // ← Cambiar a ParseUint
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	err = h.deleteUserUseCase.Execute(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{  // ← 404 más apropiado
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario eliminado correctamente",
	})
}