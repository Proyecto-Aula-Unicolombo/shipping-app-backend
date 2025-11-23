package routers

import (
	"database/sql"

	application "shipping-app/internal/app/application/orders"
	"shipping-app/internal/app/infrastructure/adapters"
	handler "shipping-app/internal/app/infrastructure/fiber/handlers/orders"

	"github.com/gofiber/fiber/v3"
)

func SetOrderRouter(apiv1 fiber.Router, db *sql.DB) {
	// Repositorios
	orderRepo := adapters.NewOrderRepositoryPostgres(db)
	driverRepo := adapters.NewDriverRepositoryAdapter(db)
	vehicleRepo := adapters.NewVehicleRepositoryPostgres(db)
	packageRepo := adapters.NewPackageRepositoryPostgres(db)
	txProvider := adapters.NewSQLTxProvider(db)
	addressPackage := adapters.NewAddressPackageRepositoryPostgres(db)
	comercialInformation := adapters.NewComercialInformationRepositoryPostgres(db)
	receiverRepo := adapters.NewReceiverRepositoryPostgres(db)

	// Casos de uso
	createOrderUC := application.NewCreateOrderUseCase(
		orderRepo,
		driverRepo,
		vehicleRepo,
		packageRepo,
		txProvider,
	)
	listOrdersUC := application.NewListOrdersUseCase(orderRepo)
	getOrderUC := application.NewGetOrderUseCase(orderRepo, packageRepo, addressPackage, comercialInformation, receiverRepo, driverRepo, vehicleRepo)
	assignOrderUC := application.NewAssignOrderUseCase(orderRepo, driverRepo, vehicleRepo, txProvider)
	updateStatusUC := application.NewUpdateOrderStatusUseCase(orderRepo)
	deleteOrderUC := application.NewDeleteOrderUseCase(orderRepo, packageRepo, txProvider)
	listByDriverUC := application.NewListOrdersByDriverUseCase(orderRepo)
	listUnassignedUC := application.NewListOrdersUnassignedUseCase(orderRepo)

	// Handler
	orderHandler := handler.NewOrderHandler(
		createOrderUC,
		listOrdersUC,
		getOrderUC,
		assignOrderUC,
		updateStatusUC,
		deleteOrderUC,
		listByDriverUC,
		listUnassignedUC,
	)

	// Rutas
	apiv1.Post("/orders", orderHandler.CreateOrder)
	apiv1.Get("/orders", orderHandler.ListOrders)
	apiv1.Get("/orders/unassigned", orderHandler.ListOrdersUnassigned)
	apiv1.Get("/orders/:id", orderHandler.GetOrder)
	apiv1.Put("/orders/:id/assign", orderHandler.AssignOrder)
	apiv1.Put("/orders/:id/status", orderHandler.UpdateStatus)
	apiv1.Delete("/orders/:id", orderHandler.DeleteOrder)
	apiv1.Get("/orders/driver/:driverId", orderHandler.ListOrdersByDriver)
}
