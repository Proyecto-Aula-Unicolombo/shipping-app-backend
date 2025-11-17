package routers

import (
	"database/sql"

	authApp "shipping-app/internal/app/application/auth"
	"shipping-app/internal/app/infrastructure/adapters"
	authHandler "shipping-app/internal/app/infrastructure/fiber/handlers/auth"
	"shipping-app/internal/externalServices/auth"

	"github.com/gofiber/fiber/v3"
)

func SetAuthRouter(app fiber.Router, db *sql.DB, jwtService *auth.JWTService) {
	userRepo := adapters.NewUserRepositoryPostgres(db)
	driverRepo := adapters.NewDriverRepositoryAdapter(db)

	loginUseCase := authApp.NewLoginUseCase(userRepo, driverRepo, jwtService)

	handler := authHandler.NewAuthHandler(loginUseCase)

	// Ruta pública de autenticación
	app.Post("/auth/login", handler.Login)
}
