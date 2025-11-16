package routers

import (
	"database/sql"
	application "shipping-app/internal/app/application/users/drivers"
	"shipping-app/internal/app/infrastructure/adapters"
	drivershandler "shipping-app/internal/app/infrastructure/fiber/handlers/users/driversHandler"

	"github.com/gofiber/fiber/v3"
)

func SetDriverRouter(apiv1 fiber.Router, db *sql.DB) {
	txRepo := adapters.NewSQLTxProvider(db)
	driverRepo := adapters.NewDriverRepositoryAdapter(db)
	userRepo := adapters.NewUserRepositoryPostgres(db)

	createDriverUseCase := application.NewCreateDriverUseCase(userRepo, driverRepo, txRepo)

	drivershandler := drivershandler.NewHandlerDrivers(createDriverUseCase)

	apiv1.Post("/drivers", drivershandler.CreateDriver)
}
