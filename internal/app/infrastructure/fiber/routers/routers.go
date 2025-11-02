package routers

import (
	"database/sql"
	// "os"
	// "shipping-app/internal/gateway/auth"
	"shipping-app/internal/gateway/services"
	// "shipping-app/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func SetupRouters(app *fiber.App, db *sql.DB) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization"},
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
	// jwtSecret := os.Getenv("JWT_SECRET")
	// jwtService := auth.NewJWTService(jwtSecret)
	apiKeyService := services.NewAPIKeyService(db)

	gateway := app.Group("/gateway")
	SetGatewayRouter(gateway, db, apiKeyService)

	apiv1 := app.Group("/api/v1")
	// apiv1.Use(middleware.JWTAuth(jwtService))
	SetUserRouter(apiv1, db)
	SetPackageRouter(apiv1, db)

}
