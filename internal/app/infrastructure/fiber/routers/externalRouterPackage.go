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
	senderHandler := externalHandler.NewSenderHandler(apiKeyService)
	external.Post("/senders/register", senderHandler.RegisterSender)

	api := external.Group("/api/v1")
	api.Use(middleware.APIKeyAuth(apiKeyService))

	addressPackage := adapters.NewAddressPackageRepositoryPostgres(db)
	comercialInformation := adapters.NewComercialInformationRepositoryPostgres(db)
	senderRepo := adapters.NewSenderRepositoryPostgres(db)
	receiverRepo := adapters.NewReceiverRepositoryPostgres(db)
	txProviderRepo := adapters.NewSQLTxProvider(db)
	domainSvc := services.NewValidatePackageService()
	repoPackage := adapters.NewPackageRepositoryPostgres(db)
	infoDeliveryRepo := adapters.NewInformationDeliveryRepositoryPostgres(db)

	createPackageUseCase := usepackages.NewCreatePackageUseCase(
		txProviderRepo, repoPackage, addressPackage, comercialInformation,
		senderRepo, receiverRepo, domainSvc,
	)
	consultPackageUseCase := usepackages.NewConsultPackageUseCase(
		repoPackage, addressPackage, comercialInformation,
		senderRepo, receiverRepo, infoDeliveryRepo,
	)
	cancelPackageUseCase := usepackages.NewCancellPackageUseCase(
		repoPackage, comercialInformation, txProviderRepo,
	)

	listPackagesUseCase := usepackages.NewListPackagesUseCase(
		repoPackage, addressPackage, comercialInformation,
		senderRepo, receiverRepo,
	)
	packageHandler := handlerpackages.NewPackageHandler(createPackageUseCase, cancelPackageUseCase, consultPackageUseCase, listPackagesUseCase, nil)

	api.Post("/packages", packageHandler.CreatePackage)
	api.Get("/packages/number/:numPackage", packageHandler.ConsultPackageByNumPackage)
	api.Delete("/packages/:numPackage", packageHandler.DeletePackage)
	api.Get("/packages", packageHandler.ListPackages)
}
