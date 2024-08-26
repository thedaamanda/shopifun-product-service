package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func UserIdHeader(c *fiber.Ctx) error {
	userId := c.Get("X-USER-ID")
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	if userId == "" {
		log.Error().Msg("middleware::UserIdHeader - Unauthorized [Header not set]")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	c.Locals("user_id", userId)

	return c.Next()
}
