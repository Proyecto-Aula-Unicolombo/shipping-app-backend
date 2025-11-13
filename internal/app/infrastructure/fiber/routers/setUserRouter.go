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
	driverRepo := adapters.NewDriverRepositoryAdapter(db)
	txProvider := adapters.NewSQLTxProvider(db)
	createUserUseCase := application.NewCreateUserUseCase(repoUser, driverRepo, txProvider)
	handlerUser := handler.NewHandlerUser(createUserUseCase)
	apiv1.Post("/users", handlerUser.CreateUser)
}
