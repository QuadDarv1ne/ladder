package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func Health(c *fiber.Ctx) error {
	return c.JSON(healthResponse{
		Status:  "ok",
		Version: version,
	})
}
