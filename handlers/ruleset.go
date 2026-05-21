package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

var exposeRuleset = os.Getenv("EXPOSE_RULESET") != "false"

func Ruleset(c *fiber.Ctx) error {
	if !exposeRuleset {
		return c.Status(fiber.StatusForbidden).SendString("Rules Disabled")
	}

	body, err := yaml.Marshal(GetRuleset())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString(string(body))
}
