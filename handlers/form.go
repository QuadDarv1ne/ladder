package handlers

import (
	_ "embed"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

//go:embed form.html
var formHtml string

func init() {
	if formPath := os.Getenv("FORM_PATH"); formPath != "" {
		dat, err := os.ReadFile(formPath)
		if err != nil {
			log.Println("ERROR: unable to load custom form", err)
		} else {
			formHtml = string(dat)
		}
	}
}

func Form(c *fiber.Ctx) error {
	if os.Getenv("DISABLE_FORM") == "true" {
		c.Set("Content-Type", "text/html")
		return c.Status(fiber.StatusNotFound).SendString("Form Disabled")
	}
	c.Set("Content-Type", "text/html")
	return c.SendString(formHtml)
}
