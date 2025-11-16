package routers

import (
	"database/sql"

	deliveryApp "shipping-app/internal/app/application/delivery"
	"shipping-app/internal/app/infrastructure/adapters"
	deliveryHandler "shipping-app/internal/app/infrastructure/fiber/handlers/delivery"

	"github.com/gofiber/fiber/v3"
)

func SetDeliveryRouter(apiv1 fiber.Router, db *sql.DB) {
	// Repositorios
	infoDeliveryRepo := adapters.NewInformationDeliveryRepositoryPostgres(db)
	packageRepo := adapters.NewPackageRepositoryPostgres(db)
	txProvider := adapters.NewSQLTxProvider(db)

	// Casos de uso
	reportDeliveryUC := deliveryApp.NewReportDeliveryUseCase(
		infoDeliveryRepo,
		packageRepo,
		txProvider,
	)

	reportIncidentUC := deliveryApp.NewReportIncidentUseCase(
		infoDeliveryRepo,
		packageRepo,
		txProvider,
	)

	getPackageReportUC := deliveryApp.NewGetPackageReportUseCase(
		infoDeliveryRepo,
		packageRepo,
	)

	// Handler
	handler := deliveryHandler.NewDeliveryHandler(
		reportDeliveryUC,
		reportIncidentUC,
		getPackageReportUC,
	)

	// Rutas
	delivery := apiv1.Group("/delivery")
	{
		// Reportar entrega exitosa
		delivery.Post("/report", handler.ReportDelivery)

		// Reportar incidente
		delivery.Post("/incident", handler.ReportIncident)

		// Obtener reporte de paquete por ID de paquete
		delivery.Get("/package/:packageId/report", handler.GetPackageReportByPackageID)

		// Obtener reporte por ID de reporte
		delivery.Get("/report/:reportId", handler.GetPackageReportByReportID)
	}
}
