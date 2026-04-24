package routers

import (
	"database/sql"
	trackingApp "shipping-app/internal/app/application/tracking"
	"shipping-app/internal/app/infrastructure/adapters"
	trackingHandler "shipping-app/internal/app/infrastructure/fiber/handlers/tracking"

	"github.com/gofiber/fiber/v3"
)

// SetPublicTrackingRouter - Rutas públicas de tracking para clientes (sin autenticación)
func SetPublicTrackingRouter(apiv1 fiber.Router, db *sql.DB) {
	packageRepo := adapters.NewPackageRepositoryPostgres(db)
	trackRepo := adapters.NewTrackRepositoryPostgres(db)
	orderRepo := adapters.NewOrderRepositoryPostgres(db)
	addressRepo := adapters.NewAddressPackageRepositoryPostgres(db)
	reciverRepo := adapters.NewReceiverRepositoryPostgres(db)

	trackPackageUC := trackingApp.NewTrackPackageUseCase(
		packageRepo,
		trackRepo,
		orderRepo,
		addressRepo,
		reciverRepo,
	)

	handler := trackingHandler.NewTrackingHandler(
		trackPackageUC,
		nil,
		nil,
		nil,
	)

	apiv1.Get("/public/tracking/:num_package", handler.TrackPackageByNumber)
}
