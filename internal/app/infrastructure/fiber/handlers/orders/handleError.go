package orders

import (
	"errors"
	"log"

	"shipping-app/internal/app/application/orders"
	"shipping-app/internal/app/domain/ports/repository"

	"github.com/gofiber/fiber/v3"
)

func (h *OrderHandler) handleErrorCreate(ctx fiber.Ctx, err error) error {
	log.Printf("Handling error: %v (type: %T)", err, err)

	switch {
	case errors.Is(err, orders.ErrInvalidOrderInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_input",
			"message": "Invalid order input data",
		})
	case errors.Is(err, orders.ErrDriverNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "driver_not_found",
			"message": "Driver not found",
		})
	case errors.Is(err, orders.ErrVehicleNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "vehicle_not_found",
			"message": "Vehicle not found",
		})
	case errors.Is(err, orders.ErrNoPackagesProvided):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_packages",
			"message": "No packages provided for order",
		})
	case errors.Is(err, orders.ErrPackageNotAvailable):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "package_not_available",
			"message": "One or more packages are not available",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_server_error",
			"message": "Could not create order",
		})
	}
}

func (h *OrderHandler) handleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repository.ErrOrderNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "order_not_found",
			"message": "Order not found",
		})
	case errors.Is(err, orders.ErrOrderCannotBeReassigned):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "order_cannot_be_reassigned",
			"message": "Order cannot be reassigned",
		})
	case errors.Is(err, orders.ErrInvalidStatus):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_status",
			"message": "Invalid order status",
		})
	case errors.Is(err, orders.ErrCannotDeleteOrder):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "cannot_delete_order",
			"message": "Cannot delete order with status other than Pendiente",
		})
	case errors.Is(err, orders.ErrDriverNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "driver_not_found",
			"message": "Driver not found",
		})
	case errors.Is(err, orders.ErrVehicleNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "vehicle_not_found",
			"message": "Vehicle not found",
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_server_error",
			"message": "An unexpected error occurred",
		})
	}
}
