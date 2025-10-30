package routers

import (
	"database/sql"
	usepackages "shipping-app/internal/app/application/UsePackages"
	services "shipping-app/internal/app/domain/services/package"
	"shipping-app/internal/app/infrastructure/adapters"
	"shipping-app/internal/app/infrastructure/fiber/handlers/handlerpackages"

	"github.com/gofiber/fiber/v3"
)

func SetPackageRouter(apiv1 fiber.Router, db *sql.DB) {
	addressPackage := adapters.NewAddressPackageRepositoryPostgres(db)
	comercialInformation := adapters.NewComercialInformationRepositoryPostgres(db)
	senderRepo := adapters.NewSenderRepositoryPostgres(db)
	receiverRepo := adapters.NewReceiverRepositoryPostgres(db)
	statusDelivery := adapters.NewStatusDeliveryRepositoryPostgres(db)
	txProviderRepo := adapters.NewSQLTxProvider(db)
	domainSvc := services.NewValidatePackageService()
	repoPackage := adapters.NewPackageRepositoryPostgres(db)

	createPackageUseCase := usepackages.NewCreatePackageUseCase(txProviderRepo, repoPackage, addressPackage, comercialInformation, senderRepo, receiverRepo, statusDelivery, domainSvc)

	packageHandler := handlerpackages.NewPackageHandler(createPackageUseCase)

	apiv1.Post("/packages", packageHandler.CreatePackage)
}
