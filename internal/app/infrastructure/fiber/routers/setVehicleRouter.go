package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/Vehicles"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/vehicles"

	"github.com/gofiber/fiber/v3"
)

func SetVehicleRouter(apiv1 fiber.Router, db *sql.DB) {
	// Repositorio y provider
	repoVehicle := adapters.NewVehicleRepositoryPostgres(db)
	txProvider := adapters.NewSQLTxProvider(db)

	// Casos de uso
	createVehicleUC := application.NewCreateVehicleUseCase(
		repoVehicle,
		txProvider,
	)
	getVehicleUC := application.NewGetVehicle(repoVehicle)  // ← AGREGAR

	// Handler con AMBOS casos de uso
	handlerVehicle := handler.NewHandlerVehicle(
		createVehicleUC,
		getVehicleUC,  // ← AGREGAR como segundo parámetro
	)

	// Rutas
	apiv1.Post("/vehicles", handlerVehicle.CreateVehicle)
	apiv1.Get("/vehicles/:id", handlerVehicle.GetVehicle)  // ← AGREGAR
}