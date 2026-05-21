package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Raw(c *fiber.Ctx) error {
	// Get the url from the URL
	urlQuery := c.Params("*")

	queries := c.Queries()
	body, _, _, err := fetchSite(urlQuery, queries)
	if err != nil {
		log.Println("ERROR:", err)
		c.Set("Content-Type", "text/plain")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Set("Content-Type", "text/plain")
	return c.SendString(body)
}
