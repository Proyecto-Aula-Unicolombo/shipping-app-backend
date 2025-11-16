package routers

import (
	"database/sql"

	trackingApp "shipping-app/internal/app/application/tracking"
	"shipping-app/internal/app/infrastructure/adapters"
	trackingHandler "shipping-app/internal/app/infrastructure/fiber/handlers/tracking"

	"github.com/gofiber/fiber/v3"
)

func SetTrackingRouter(apiv1 fiber.Router, db *sql.DB) {
	// Repositorios
	packageRepo := adapters.NewPackageRepositoryPostgres(db)
	trackRepo := adapters.NewTrackRepositoryPostgres(db)
	orderRepo := adapters.NewOrderRepositoryPostgres(db)
	stopRepo := adapters.NewDeliveryStopRepositoryPostgres(db)
	txProvider := adapters.NewSQLTxProvider(db)

	// Casos de uso
	trackPackageUC := trackingApp.NewTrackPackageUseCase(
		packageRepo,
		trackRepo,
		orderRepo,
	)

	registerStopUC := trackingApp.NewRegisterStopUseCase(
		stopRepo,
		trackRepo,
		orderRepo,
		txProvider,
	)

	// Handler
	handler := trackingHandler.NewTrackingHandler(
		trackPackageUC,
		registerStopUC,
	)

	// Rutas públicas de tracking (para destinatarios)
	tracking := apiv1.Group("/tracking")
	{
		// Rastrear paquete por número (público para destinatario)
		tracking.Get("/package", handler.TrackPackageByNumber)

		// Rastrear paquete por ID (interno)
		tracking.Get("/package/:packageId", handler.TrackPackageByID)
	}

	// Rutas privadas para conductores
	stops := apiv1.Group("/stops")
	{
		// Registrar parada durante entrega
		stops.Post("/register", handler.RegisterStop)
	}
}
