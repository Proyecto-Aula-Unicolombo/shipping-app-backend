package gateway

import (
	"shipping-app/internal/gateway/services"

	"github.com/gofiber/fiber/v3"
)

type SenderHandler struct {
	apiKeyService *services.APIKeyService
}

func NewSenderHandler(apiKeyService *services.APIKeyService) *SenderHandler {
	return &SenderHandler{
		apiKeyService: apiKeyService,
	}
}

type RegisterSenderRequest struct {
	Name        string `json:"name"`
	Document    string `json:"document"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (h *SenderHandler) RegisterSender(c fiber.Ctx) error {
	var req RegisterSenderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" || req.Document == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "missing_fields",
			"message": "Name, document and email are required",
		})
	}

	sender, apiKey, err := h.apiKeyService.CreateSenderWithAPIKey(
		req.Name,
		req.Document,
		req.Address,
		req.PhoneNumber,
		req.Email,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "registration_failed",
			"message": "Could not register sender",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sender_id": sender.ID,
		"name":      sender.Name,
		"document":  sender.Document,
		"email":     sender.Email,
		"api_key":   apiKey,
		"message":   "Sender registered successfully. Save your API key, it won't be shown again.",
	})
}
