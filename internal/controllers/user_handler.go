package controllers

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type userServiceHandler struct {
	UserService entities.UserService
}

func NewUserServiceHandler(app *fiber.App, u entities.UserService) {
	handler := &userServiceHandler{UserService: u}
	app.Post("/user/register", handler.RegisterHandler)
	app.Post("/user/login", handler.LoginHandler)
	app.Get("/refresh", handler.RefreshHandler)
}

func (u *userServiceHandler) RegisterHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "invalid username",
		})
	}
	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"note": "invalid password",
		})
	}
	err := u.UserService.Register(c.Context(), username, password)
	if err != nil {
		return errorHandling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (u *userServiceHandler) LoginHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	access_token, refresh_token, err := u.UserService.Login(c.Context(), username, password)
	if err != nil {
		return errorHandling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

func (u *userServiceHandler) RefreshHandler(c *fiber.Ctx) error {
	tokenstring := c.Get("Authorization")[7:]
	access_token, refresh_token, err := u.UserService.Refresh(tokenstring)
	if err != nil {
		return errorHandling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}
