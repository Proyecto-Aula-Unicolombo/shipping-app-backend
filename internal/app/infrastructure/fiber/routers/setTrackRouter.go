package routers

import (
	"database/sql"
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters"
	tracksHandlers "shipping-app/internal/app/infrastructure/fiber/handlers/tracksHabdkers"

	"github.com/gofiber/fiber/v3"
)

func SetTrackRouter(apiv1 fiber.Router, db *sql.DB) {
	trackRepository := adapters.NewTrackingRepositoryAdapter(db)

	registerTrack := tracks.NewTrackRegisterUseCase(trackRepository)

	handlerTrack := tracksHandlers.NewTrackHandler(registerTrack)

	apiv1.Post("/tracks", handlerTrack.RegisterTrack)
}
