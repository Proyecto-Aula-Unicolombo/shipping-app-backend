package routers

import (
	"database/sql"

	usepackages "shipping-app/internal/app/application/UsePackages"
	services "shipping-app/internal/app/domain/services/package"
	"shipping-app/internal/app/infrastructure/adapters"
	externalHandler "shipping-app/internal/app/infrastructure/fiber/handlers/externalHandler"
	"shipping-app/internal/app/infrastructure/fiber/handlers/handlerpackages"
	api_key "shipping-app/internal/externalServices/services"
	"shipping-app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetExternalRouter(external fiber.Router, db *sql.DB, apiKeyService *api_key.APIKeyService) {
	// Rutas públicas
	senderHandler := externalHandler.NewSenderHandler(apiKeyService)
	external.Post("/senders/register", senderHandler.RegisterSender)

	// Rutas protegidas con API Key
	api := external.Group("/api/v1")
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
