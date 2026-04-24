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
	addressRepo := adapters.NewAddressPackageRepositoryPostgres(db)
	reciverRepo := adapters.NewReceiverRepositoryPostgres(db)

	// Casos de uso
	trackPackageUC := trackingApp.NewTrackPackageUseCase(
		packageRepo,
		trackRepo,
		orderRepo,
		addressRepo,
		reciverRepo,
	)

	registerStopUC := trackingApp.NewRegisterStopUseCase(
		stopRepo,
		trackRepo,
		orderRepo,
		txProvider,
	)

	listIncidentsUC := trackingApp.NewListIncidentsUseCase(
		stopRepo,
		orderRepo,
	)

	listActiveOrdersUC := trackingApp.NewListActiveOrdersUseCase(
		orderRepo,
		trackRepo,
	)

	// Handler
	handler := trackingHandler.NewTrackingHandler(
		trackPackageUC,
		registerStopUC,
		listIncidentsUC,
		listActiveOrdersUC,
	)

	tracking := apiv1.Group("/tracking")
	{

		// Rastrear paquete por ID (interno)
		tracking.Get("/package/:packageId", handler.TrackPackageByID)
	}

	// Rutas privadas para conductores
	stops := apiv1.Group("/stops")
	{
		// Registrar parada durante entrega
		stops.Post("/register", handler.RegisterStop)
	}

	// Rutas de incidentes (acceso para admins/coordinadores)
	incidents := apiv1.Group("/incidents")
	{
		// Listar incidentes con filtros
		incidents.Get("/", handler.ListIncidents)
	}

	// Rutas de tracking para admin/coordinador
	adminTracking := apiv1.Group("/admin/tracking")
	{
		// Listar todas las órdenes activas con ubicación
		adminTracking.Get("/active-orders", handler.ListActiveOrders)
	}
}
