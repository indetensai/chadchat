package auth

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

func New(repo entities.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenstring := c.Get("Authorization")[7:]
		if tokenstring == "" {
			c.Locals("token", nil)
			return c.Next()
		}
		user_id, username, err := repo.Validation(tokenstring)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("invalid token")
		}
		c.Locals("user_id", user_id)
		c.Locals("username", username)
		return c.Next()
	}
}
