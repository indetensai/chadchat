package controllers

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

func error_handling(c *fiber.Ctx, err error) error {
	switch err {
	case entities.ErrDuplicate:
		return c.SendStatus(fiber.StatusConflict)
	case entities.ErrNotAuthorized:
		return c.SendStatus(fiber.StatusForbidden)
	case entities.ErrNotFound:
		return c.SendStatus(fiber.StatusNotFound)
	case entities.ErrEmptySession:
		return c.SendStatus(fiber.StatusUnauthorized)
	case entities.ErrInvalidCredentials:
		return c.SendStatus(fiber.StatusUnauthorized)
	case nil:
		return nil
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}
