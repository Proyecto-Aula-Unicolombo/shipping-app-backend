package routers

import (
	"database/sql"
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// SetPublicTrackingRouter - Rutas públicas de tracking para clientes (sin autenticación)
func SetPublicTrackingRouter(apiv1 fiber.Router, db *sql.DB) {
	trackRepository := adapters.NewTrackRepositoryPostgres(db)
	orderRepository := adapters.NewOrderRepositoryPostgres(db)

	// Caso de uso para obtener tracking de una orden
	getOrderTracks := tracks.NewGetOrderTracksUseCase(trackRepository, orderRepository)

	// Handler dedicado para tracking público
	handler := &PublicTrackingHandler{
		getOrderTracksUC: getOrderTracks,
	}

	// Ruta pública: /api/v1/public/tracking/:orderNumber
	apiv1.Get("/public/tracking/:orderNumber", handler.GetOrderTracking)
}

type PublicTrackingHandler struct {
	getOrderTracksUC *tracks.GetOrderTracksUseCase
}

// GetOrderTracking - Endpoint público para rastreo de pedidos por número de orden
func (h *PublicTrackingHandler) GetOrderTracking(ctx fiber.Ctx) error {
	orderNumber := ctx.Params("orderNumber")
	if orderNumber == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_order",
			"message": "Order number is required",
		})
	}

	// Buscar orden por número (asumiendo que el número es el ID)
	// En producción, podrías tener un campo NumOrder diferente al ID
	orderIDUint64, err := strconv.ParseUint(orderNumber, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_order",
			"message": "Order number must be a valid number",
		})
	}
	orderID := uint(orderIDUint64)

	// Obtener tracking de la orden
	input := tracks.GetOrderTracksInput{
		OrderID: orderID,
		Limit:   nil, // Obtener todos los tracks
	}

	output, err := h.getOrderTracksUC.Execute(ctx.Context(), input)
	if err != nil {
		switch err {
		case tracks.ErrOrderNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "not_found",
				"message": "Order not found. Please verify the order number.",
			})
		case tracks.ErrNoTracksFound:
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"data": fiber.Map{
					"order_id": orderID,
					"status":   output.Status,
					"tracks":   []interface{}{},
					"message":  "No tracking data available yet",
				},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "internal_error",
				"message": "Failed to retrieve tracking data",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"data": output,
	})
}
