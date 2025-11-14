package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/users"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/users"

	"github.com/gofiber/fiber/v3"
)

func SetUserRouter(apiv1 fiber.Router, db *sql.DB) {
	// Repositorios
	repoUser := adapters.NewUserRepositoryPostgres(db)
	driverRepo := adapters.NewDriverRepositoryAdapter(db)
	txProvider := adapters.NewSQLTxProvider(db)

	// Casos de uso
	createUserUseCase := application.NewCreateUserUseCase(repoUser, driverRepo, txProvider)
	getUserUseCase := application.NewGetUser(repoUser, driverRepo)
	deleteUserUseCase := application.NewDeleteUserUseCase(repoUser)
	listUsersUseCase := application.NewListUsers(repoUser)            // Tuyo (sin paginación)
	listUsersPaginatedUC := application.NewListUsersUseCase(repoUser) // Del compañero (con paginación)
	updateUserUseCase := application.NewUpdateUserUseCase(repoUser)

	// Handler con TODOS los casos de uso
	userHandler := handler.NewHandlerUser(
		createUserUseCase,
		getUserUseCase,
		deleteUserUseCase,
		listUsersUseCase,
		listUsersPaginatedUC,
		updateUserUseCase,
	)

	// Rutas
	apiv1.Post("/users", userHandler.CreateUser)
	apiv1.Get("/users/:id", userHandler.GetUser)
	apiv1.Get("/users", userHandler.ListUsersPaginated)  // Con paginación (del compañero)
	apiv1.Get("/users/all", userHandler.ListUsersSimple) // Sin paginación (tuyo)
	apiv1.Put("/users/:id", userHandler.UpdateUser)
	apiv1.Delete("/users/:id", userHandler.DeleteUser)
}
