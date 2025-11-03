package routers

import (
	"database/sql"
	usepackages "shipping-app/internal/app/application/UsePackages"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/app/infrastructure/fiber/handlers/handlerpackages"

	"github.com/gofiber/fiber/v3"
)

func SetPackageRouter(api fiber.Router, db *sql.DB) {
	addressPackage := adapters.NewAddressPackageRepositoryPostgres(db)
	comercialInformation := adapters.NewComercialInformationRepositoryPostgres(db)
	senderRepo := adapters.NewSenderRepositoryPostgres(db)
	receiverRepo := adapters.NewReceiverRepositoryPostgres(db)
	statusDelivery := adapters.NewStatusDeliveryRepositoryPostgres(db)
	// txProviderRepo := adapters.NewSQLTxProvider(db)
	repoPackage := adapters.NewPackageRepositoryPostgres(db)

	consultPackageUseCase := usepackages.NewConsultPackageUseCase(repoPackage, addressPackage, comercialInformation, senderRepo, receiverRepo, statusDelivery)
	listPackagesUseCase := usepackages.NewListPackagesUseCase(
		repoPackage, addressPackage, comercialInformation,
		senderRepo, receiverRepo, statusDelivery,
	)
	packageHandler := handlerpackages.NewPackageHandler(nil, nil, consultPackageUseCase, listPackagesUseCase)

	api.Get("/packages/:id", packageHandler.ConsultPackageByID)
	api.Get("/packages", packageHandler.ListPackages)
}
