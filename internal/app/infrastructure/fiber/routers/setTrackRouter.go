package routers

import (
	"database/sql"
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	tracksHandlers "shipping-app/internal/app/infrastructure/fiber/handlers/tracksHabdkers"

	"github.com/gofiber/fiber/v3"
)

func SetTrackRouter(apiv1 fiber.Router, db *sql.DB, hub *ws.Hub) {
	trackingRepository := adapters.NewTrackingRepositoryAdapter(db)
	trackRepository := adapters.NewTrackRepositoryPostgres(db)
	orderRepository := adapters.NewOrderRepositoryPostgres(db)
	driverRepository := adapters.NewDriverRepositoryAdapter(db)

	// Casos de uso
	registerTrack := tracks.NewTrackRegisterUseCase(trackingRepository)
	getOrderTracks := tracks.NewGetOrderTracksUseCase(trackRepository, orderRepository)
	getActiveDriversLocations := tracks.NewGetActiveDriversLocationsUseCase(
		trackRepository,
		orderRepository,
		driverRepository,
	)

	handlerTrack := tracksHandlers.NewTrackHandler(
		registerTrack,
		getOrderTracks,
		getActiveDriversLocations,
		hub,
	)

	// Rutas
	apiv1.Post("/tracks", handlerTrack.RegisterTrack)
	apiv1.Get("/tracks/order/:orderId", handlerTrack.GetOrderTracks)
	apiv1.Get("/tracks/active-drivers", handlerTrack.GetActiveDriversLocations)
}
