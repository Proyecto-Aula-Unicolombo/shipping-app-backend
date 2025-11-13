package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/users"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/users"

	"github.com/gofiber/fiber/v3"
)

func SetUserRouter(apiv1 fiber.Router, db *sql.DB) {

	repoUser := adapters.NewUserRepositoryPostgres(db)

	createUserUseCase := application.NewCreateUserUseCase(repoUser)
	getUserUseCase := application.NewGetUser(repoUser)
	deleteUserUseCase := application.NewDeleteUserUseCase(repoUser)
	listUsersUseCase := application.NewListUsers(repoUser)
	updateUserUseCase := application.NewUpdateUserUseCase(repoUser)

	userHandler := handler.NewHandlerUser(
		createUserUseCase,
		getUserUseCase,
		deleteUserUseCase,
		listUsersUseCase,
		updateUserUseCase,
	)

	apiv1.Post("/users", userHandler.CreateUser)
	apiv1.Get("/users/:id", userHandler.GetUser)
	apiv1.Get("/users", userHandler.ListUsers)
	apiv1.Delete("/users/:id", userHandler.DeleteUser)
	apiv1.Put("/users/:id", userHandler.UpdateUser)
}
