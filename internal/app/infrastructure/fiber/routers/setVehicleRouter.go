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
		deleteVehicleUC := application.NewDeleteVehicleUseCase(repoVehicle) 
	getVehicleUC := application.NewGetVehicle(repoVehicle)  
	
	listVehiclesUC := application.NewListVehicles(repoVehicle) // ← AGREGAR

	// Handler con AMBOS casos de uso
	handlerVehicle := handler.NewHandlerVehicle(
		createVehicleUC,
		getVehicleUC,  // ← AGREGAR como segundo parámetro
		deleteVehicleUC,
		listVehiclesUC,
	)

	// Rutas
	apiv1.Post("/vehicles", handlerVehicle.CreateVehicle)
	apiv1.Get("/vehicles/:id", handlerVehicle.GetVehicle)
		apiv1.Delete("/vehicles/:id", handlerVehicle.DeleteVehicle) 
		apiv1.Get("/vehicles", handlerVehicle.ListVehiclesSimple)  // ← AGREGAR
}