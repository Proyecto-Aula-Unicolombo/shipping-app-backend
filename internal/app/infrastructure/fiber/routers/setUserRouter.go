package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/users"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/users"
	"shipping-app/internal/middleware"
	"shipping-app/internal/externalServices/auth"

	"github.com/gofiber/fiber/v3"
)

func SetUserRouter(apiv1 fiber.Router, db *sql.DB, jwtService *auth.JWTService) {
	repoUser := adapters.NewUserRepositoryPostgres(db)
	driverRepo := adapters.NewDriverRepositoryAdapter(db)
	txProvider := adapters.NewSQLTxProvider(db)

	createUserUseCase := application.NewCreateUserUseCase(repoUser, driverRepo, txProvider)
	getUserUseCase := application.NewGetUser(repoUser, driverRepo)
	deleteUserUseCase := application.NewDeleteUserUseCase(repoUser, driverRepo, txProvider)
	listUsersUseCase := application.NewListUsers(repoUser)
	listUsersPaginatedUC := application.NewListUsersUseCase(repoUser)
	updateUserUseCase := application.NewUpdateUserUseCase(repoUser, driverRepo, txProvider)

	userHandler := handler.NewHandlerUser(
		createUserUseCase,
		getUserUseCase,
		deleteUserUseCase,
		listUsersUseCase,
		listUsersPaginatedUC,
		updateUserUseCase,
	)

	apiv1.Post("/users", userHandler.CreateUser)
	protected := apiv1.Group("", middleware.JWTAuth(jwtService))
	protected.Get("/users/:id", userHandler.GetUser)
	protected.Get("/users", userHandler.ListUsersPaginated)
	protected.Put("/users/:id", userHandler.UpdateUser)
	protected.Delete("/users/:id", userHandler.DeleteUser)
}
