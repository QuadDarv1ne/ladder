package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func computeShareToken(secret, targetURL string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(targetURL))
	return hex.EncodeToString(mac.Sum(nil))[:32]
}

func ShareLink(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Status(fiber.StatusServiceUnavailable).SendString("share feature is disabled (SHARE_SECRET not set)")
		}

		targetURL := c.Query("url")
		if targetURL == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing 'url' query parameter")
		}

		token := computeShareToken(secret, targetURL)
		shareURL := fmt.Sprintf("%s://%s/share/%s/%s", c.Protocol(), c.Hostname(), token, targetURL)

		return c.JSON(fiber.Map{
			"url": shareURL,
		})
	}
}

func ProxySiteViaShare(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Status(fiber.StatusServiceUnavailable).SendString("share feature is disabled (SHARE_SECRET not set)")
		}

		targetURL, err := extractUrl(c)
		if err != nil || targetURL == "" {
			log.Println("ERROR in share redirect URL extraction:", err)
			return c.Status(fiber.StatusBadRequest).SendString("invalid URL")
		}

		token := computeShareToken(secret, targetURL)
		return c.Redirect("/share/"+token+"/"+targetURL, fiber.StatusFound)
	}
}

func ShareProxy(secret, rulesetPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}

		providedToken := c.Params("token")

		targetURL, err := extractUrl(c)
		if err != nil {
			log.Println("ERROR in share URL extraction:", err)
			return c.Status(fiber.StatusBadRequest).SendString("invalid URL")
		}

		if targetURL == "" {
			return c.Status(fiber.StatusBadRequest).SendString("missing target URL")
		}

		expectedToken := computeShareToken(secret, targetURL)

		if !hmac.Equal([]byte(providedToken), []byte(expectedToken)) {
			return c.Status(fiber.StatusForbidden).SendString("invalid share token")
		}

		queries := c.Queries()
		body, _, resp, err := fetchSite(targetURL, queries)
		if err != nil {
			log.Println("ERROR in share proxy:", err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		c.Set("Content-Security-Policy", resp.Header.Get("Content-Security-Policy"))

		return c.SendString(body)
	}
}
