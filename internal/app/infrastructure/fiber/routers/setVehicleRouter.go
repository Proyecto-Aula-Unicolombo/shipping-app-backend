package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/vehicles"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/vehicles"

	"github.com/gofiber/fiber/v3"
)

func SetVehicleRouter(apiv1 fiber.Router, db *sql.DB) {

	repoVehicle := adapters.NewVehicleRepositoryPostgres(db)

	createVehicleUC := application.NewCreateVehicleUseCase(
		repoVehicle,
	)
	deleteVehicleUC := application.NewDeleteVehicleUseCase(repoVehicle)
	getVehicleUC := application.NewGetVehicle(repoVehicle)

	listVehiclesUC := application.NewListVehicles(repoVehicle)
	updateVehicleUC := application.NewUpdateVehicleUseCase(repoVehicle)
	listUnassignedVehiclesUC := application.NewListUnassignedVehiclesUseCase(repoVehicle)

	handlerVehicle := handler.NewHandlerVehicle(
		createVehicleUC,
		getVehicleUC,
		deleteVehicleUC,
		listVehiclesUC,
		updateVehicleUC,
		listUnassignedVehiclesUC,
	)

	apiv1.Post("/vehicles", handlerVehicle.CreateVehicle)
	apiv1.Get("/vehicles", handlerVehicle.ListVehiclesSimple)
	apiv1.Get("/vehicles/unassigned", handlerVehicle.ListUnassignedVehicles)
	apiv1.Get("/vehicles/:id", handlerVehicle.GetVehicle)
	apiv1.Delete("/vehicles/:id", handlerVehicle.DeleteVehicle)
	apiv1.Put("/vehicles/:id", handlerVehicle.UpdateVehicle)
}
