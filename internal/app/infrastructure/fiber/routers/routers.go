package routers

import (
	"database/sql"
	"os"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	"shipping-app/internal/externalServices/auth"
	"shipping-app/internal/externalServices/services"
	"shipping-app/internal/middleware"

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

	// Configurar servicios de autenticación
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // Fallback para desarrollo
	}
	jwtService := auth.NewJWTService(jwtSecret)
	apiKeyService := services.NewAPIKeyService(db)

	// Rutas externas (API Gateway con API Key)
	external := app.Group("/external")
	SetExternalRouter(external, db, apiKeyService)

	// Rutas públicas (sin autenticación)
	apiv1 := app.Group("/api/v1")
	SetAuthRouter(apiv1, db, jwtService)

	// Rutas protegidas (requieren JWT)
	protected := apiv1.Group("", middleware.JWTAuth(jwtService))
	SetUserRouter(protected, db)
	SetPackageRouter(protected, db)
	SetTrackRouter(protected, db, hub)
	SetVehicleRouter(protected, db)
	SetOrderRouter(protected, db)
	SetDriverRouter(protected, db)
	SetDeliveryRouter(protected, db)
	SetTrackingRouter(protected, db)
}
