package tracksHandlers

import (
	"log"
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type TrackHandler struct {
	registerUC                *tracks.TrackRegisterUseCase
	getOrderTracksUC          *tracks.GetOrderTracksUseCase
	getActiveDriversLocations *tracks.GetActiveDriversLocationsUseCase
	hub                       *ws.Hub
}

func NewTrackHandler(
	registerUC *tracks.TrackRegisterUseCase,
	getOrderTracksUC *tracks.GetOrderTracksUseCase,
	getActiveDriversLocations *tracks.GetActiveDriversLocationsUseCase,
	hub *ws.Hub,
) *TrackHandler {
	return &TrackHandler{
		registerUC:                registerUC,
		getOrderTracksUC:          getOrderTracksUC,
		getActiveDriversLocations: getActiveDriversLocations,
		hub:                       hub,
	}
}

type TrackRegisterRequest struct {
	OrderID   *uint    `json:"order_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
type TrackRegisterResponse struct {
	OrderID   uint    `json:"order_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
	TrackID   uint    `json:"track_id"`
}

func (h *TrackHandler) RegisterTrack(ctx fiber.Ctx) error {
	var req TrackRegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Failed to parse request body",
		})
	}
	if req.OrderID == nil || req.Latitude == nil || req.Longitude == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "OrderID, Latitude and Longitude are required",
		})
	}

	trackInput := &tracks.TrackRegisterInput{
		OrderID:   *req.OrderID,
		Latitude:  *req.Latitude,
		Longitude: *req.Longitude,
	}
	trackOutput, err := h.registerUC.Execute(ctx.Context(), trackInput)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "interval_server_error",
			"message": "could not create track",
		})
	}

	trackOutputResponse := TrackRegisterResponse{
		OrderID:   trackOutput.OrderID,
		Latitude:  trackOutput.Latitude,
		Longitude: trackOutput.Longitude,
		Timestamp: trackOutput.Timestamp.Format(time.RFC3339),
		TrackID:   trackOutput.TrackID,
	}

	var postMessage = ws.WebSocketMessage{
		Type:    "track_update",
		Payload: trackOutputResponse,
	}

	// Enviar actualización solo a admins y clientes siguiendo esta orden
	log.Printf("📡 Enviando track_update por WebSocket: order_id=%d, lat=%.4f, lng=%.4f", trackOutput.OrderID, trackOutput.Latitude, trackOutput.Longitude)
	h.hub.BroadcastToOrder(postMessage, trackOutput.OrderID)
	log.Printf("✅ Mensaje enviado via WebSocket")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Track created successfully"})
}

// GetOrderTracks - Obtener historial de ubicaciones de una orden
func (h *TrackHandler) GetOrderTracks(ctx fiber.Ctx) error {
	orderIDStr := ctx.Params("orderId")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_order_id",
			"message": "Order ID must be a valid number",
		})
	}

	// Obtener limit opcional desde query params
	var limit *int
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil && limitInt > 0 {
			limit = &limitInt
		}
	}

	input := tracks.GetOrderTracksInput{
		OrderID: uint(orderID),
		Limit:   limit,
	}

	output, err := h.getOrderTracksUC.Execute(ctx.Context(), input)
	if err != nil {
		switch err {
		case tracks.ErrOrderNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "not_found",
				"message": "Order not found",
			})
		case tracks.ErrNoTracksFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "no_tracks",
				"message": "No tracking data found for this order",
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

// GetActiveDriversLocations - Obtener ubicaciones de todos los conductores activos
func (h *TrackHandler) GetActiveDriversLocations(ctx fiber.Ctx) error {
	output, err := h.getActiveDriversLocations.Execute(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "Failed to retrieve driver locations",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": output,
	})
}
