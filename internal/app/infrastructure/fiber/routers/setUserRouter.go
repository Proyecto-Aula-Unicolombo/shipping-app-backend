package routers

import (
	"database/sql"
	application "shipping-app/internal/app/application/users"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/users"

	"shipping-app/internal/app/infrastructure/adapters"

	"github.com/gofiber/fiber/v3"
)

func SetUserRouter(apiv1 fiber.Router, db *sql.DB) {
	repoUser := adapters.NewUserRepositoryPostgres(db)
	createUserUseCase := application.NewCreateUserUseCase(repoUser)
	createUserHandler := handler.NewCreateUserHandler(createUserUseCase)
	apiv1.Post("/users", createUserHandler.CreateUserHandler)
}
