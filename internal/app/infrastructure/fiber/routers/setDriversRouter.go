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
	listDriverUseCase := application.NewListDriverUseCase(driverRepo, userRepo)
	getDriverByIdUseCase := application.NewGetByIdDriverUseCase(driverRepo)
	updateStatusDriverUseCase := application.NewUpdateStatusDriverUseCase(driverRepo)

	drivershandler := drivershandler.NewHandlerDrivers(createDriverUseCase, listDriverUseCase, getDriverByIdUseCase, updateStatusDriverUseCase)

	apiv1.Post("/drivers", drivershandler.CreateDriver)
	apiv1.Get("/drivers", drivershandler.ListDrivers)
	apiv1.Get("/drivers/:id", drivershandler.GetDriverByID)
	apiv1.Patch("/drivers/:id/status", drivershandler.UpdateStatusDriver)
}
