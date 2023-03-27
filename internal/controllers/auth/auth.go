package auth

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

func New(userService entities.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenstring := c.Get("Authorization")[7:]
		if tokenstring == "" {
			c.Locals("token", nil)
			return c.Next()
		}
		token_credentials, err := userService.Validation(tokenstring)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("invalid token")
		}
		c.Locals("user_id", token_credentials.UserID)
		c.Locals("username", token_credentials.Username)
		return c.Next()
	}
}
