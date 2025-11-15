package routers

import (
	"database/sql"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	"shipping-app/internal/externalServices/services"


	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func SetupRouters(app *fiber.App, db *sql.DB) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"Upgrade",
			"Connection",
			"Sec-WebSocket-Key",
			"Sec-WebSocket-Version",
			"Sec-WebSocket-Protocol",
		},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		},
	}))

	hub := ws.NewHub()
	go hub.Run()
	app.Get("/api/v1/ws", hub.HandleWebSocketConnection)
	apiKeyService := services.NewAPIKeyService(db)

	external := app.Group("/external")
	SetExternalRouter(external, db, apiKeyService)

	apiv1 := app.Group("/api/v1")
	SetUserRouter(apiv1, db)
	SetPackageRouter(apiv1, db)
	SetTrackRouter(apiv1, db, hub)

}
