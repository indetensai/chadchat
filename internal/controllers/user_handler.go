package controllers

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type userServiceHandler struct {
	UserService entities.UserService
}

func NewUserServiceHandler(app *fiber.App, u entities.UserService) entities.UserHandler {
	handler := &userServiceHandler{UserService: u}
	app.Post("/user/register", handler.RegisterHandler)
	app.Post("/user/login", handler.LoginHandler)
	return handler
}

func (chat *userServiceHandler) RegisterHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := chat.UserService.Register(c.Context(), username, password)
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (chat *userServiceHandler) LoginHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	access_token, refresh_token, err := chat.UserService.Login(c.Context(), username, password)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}
