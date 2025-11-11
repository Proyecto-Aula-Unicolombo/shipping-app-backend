package routers

import (
	"database/sql"
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	tracksHandlers "shipping-app/internal/app/infrastructure/fiber/handlers/tracksHabdkers"

	// "shipping-app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetTrackRouter(apiv1 fiber.Router, db *sql.DB, hub *ws.Hub) {
	trackRepository := adapters.NewTrackingRepositoryAdapter(db)

	registerTrack := tracks.NewTrackRegisterUseCase(trackRepository)

	handlerTrack := tracksHandlers.NewTrackHandler(registerTrack, hub)

	apiv1.Post("/tracks", handlerTrack.RegisterTrack)
}
