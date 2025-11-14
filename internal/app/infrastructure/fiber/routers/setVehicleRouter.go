package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/Vehicles"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/vehicles"

	"github.com/gofiber/fiber/v3"
)

func SetVehicleRouter(apiv1 fiber.Router, db *sql.DB) {
	repoVehicle := adapters.NewVehicleRepositoryPostgres(db)
	txProvider := adapters.NewSQLTxProvider(db)

	createVehicleUC := application.NewCreateVehicleUseCase(
		repoVehicle,
		txProvider,
	)

	handlerVehicle := handler.NewHandlerVehicle(createVehicleUC)

	apiv1.Post("/vehicles", handlerVehicle.CreateVehicle)
}
