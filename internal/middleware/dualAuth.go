package middleware

import (
	"strings"

	"shipping-app/internal/externalServices/auth"
	"shipping-app/internal/externalServices/services"

	"github.com/gofiber/fiber/v3"
)

// para UI interna
func JWTAuth(jwtService *auth.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := extractJWTToken(c)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "JWT token required",
			})
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "invalid_token",
				"message": "Invalid or expired JWT token",
			})
		}

		// Guardar información del usuario en el contexto
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("driver_id", claims.DriverID)

		return c.Next()
	}
}

// para API Gateway externa
func APIKeyAuth(apiKeyService *services.APIKeyService) fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := extractAPIKey(c)
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "API Key required in X-API-Key header",
			})
		}

		sender, err := apiKeyService.ValidateAPIKey(apiKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "invalid_api_key",
				"message": "Invalid or inactive API key",
			})
		}

		// Guardar información del sender en el contexto
		c.Locals("sender_id", sender.ID)
		c.Locals("sender_document", sender.Document)
		c.Locals("sender_name", sender.Name)

		return c.Next()
	}
}

func extractJWTToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func extractAPIKey(c fiber.Ctx) string {
	return c.Get("X-API-Key")
}
