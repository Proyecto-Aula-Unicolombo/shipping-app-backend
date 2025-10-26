package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func InitFiber() (*fiber.App, error) {
	appFiber := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		CaseSensitive: true,
		StrictRouting: false,
		ServerHeader:  "MyFiberApp",
	})
	return appFiber, nil
}