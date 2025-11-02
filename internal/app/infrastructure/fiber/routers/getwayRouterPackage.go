package routers

import (
	"database/sql"

	usepackages "shipping-app/internal/app/application/UsePackages"
	services "shipping-app/internal/app/domain/services/package"
	"shipping-app/internal/app/infrastructure/adapters"
	handlergateway "shipping-app/internal/app/infrastructure/fiber/handlers/gateway"
	"shipping-app/internal/app/infrastructure/fiber/handlers/handlerpackages"
	gatewayservices "shipping-app/internal/gateway/services"
	"shipping-app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetGatewayRouter(gateway fiber.Router, db *sql.DB, apiKeyService *gatewayservices.APIKeyService) {
	// Rutas públicas
	senderHandler := handlergateway.NewSenderHandler(apiKeyService)
	gateway.Post("/senders/register", senderHandler.RegisterSender)

	// Rutas protegidas con API Key
	api := gateway.Group("/api/v1")
	api.Use(middleware.APIKeyAuth(apiKeyService))

	addressPackage := adapters.NewAddressPackageRepositoryPostgres(db)
	comercialInformation := adapters.NewComercialInformationRepositoryPostgres(db)
	senderRepo := adapters.NewSenderRepositoryPostgres(db)
	receiverRepo := adapters.NewReceiverRepositoryPostgres(db)
	statusDelivery := adapters.NewStatusDeliveryRepositoryPostgres(db)
	txProviderRepo := adapters.NewSQLTxProvider(db)
	domainSvc := services.NewValidatePackageService()
	repoPackage := adapters.NewPackageRepositoryPostgres(db)

	createPackageUseCase := usepackages.NewCreatePackageUseCase(
		txProviderRepo, repoPackage, addressPackage, comercialInformation,
		senderRepo, receiverRepo, statusDelivery, domainSvc,
	)
	consultPackageUseCase := usepackages.NewConsultPackageUseCase(
		repoPackage, txProviderRepo, addressPackage, comercialInformation,
		senderRepo, receiverRepo, statusDelivery,
	)
	cancelPackageUseCase := usepackages.NewCancellPackageUseCase(
		repoPackage, comercialInformation, statusDelivery, txProviderRepo,
	)

	packageHandler := handlerpackages.NewPackageHandler(createPackageUseCase, cancelPackageUseCase, consultPackageUseCase)

	api.Post("/packages", packageHandler.CreatePackage)
	api.Get("/packages/number/:numPackage", packageHandler.ConsultPackageByNumPackage)
	api.Delete("/packages/:numPackage", packageHandler.DeletePackage)
}
