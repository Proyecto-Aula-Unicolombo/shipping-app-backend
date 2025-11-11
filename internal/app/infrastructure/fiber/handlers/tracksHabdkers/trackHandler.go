package tracksHandlers

import (
	"shipping-app/internal/app/application/tracks"
	"shipping-app/internal/app/infrastructure/adapters/ws"
	"time"

	"github.com/gofiber/fiber/v3"
)

type TrackHandler struct {
	registerUC *tracks.TrackRegisterUseCase
	hub        *ws.Hub
}

func NewTrackHandler(registerUC *tracks.TrackRegisterUseCase, hub *ws.Hub) *TrackHandler {
	return &TrackHandler{registerUC: registerUC, hub: hub}
}

type TrackRegisterRequest struct {
	OrderID   *uint    `json:"order_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
type TrackRegisterResponse struct {
	OrderID   uint      `json:"order_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	TrackID   uint      `json:"track_id"`
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
		Timestamp: trackOutput.Timestamp,
		TrackID:   trackOutput.TrackID,
	}

	var postMessage = ws.WebSocketMessage{
		Type:    "Track_created",
		Payload: trackOutputResponse,
	}

	h.hub.BroadcastJSON(postMessage, nil)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Track created successfully"})
}
