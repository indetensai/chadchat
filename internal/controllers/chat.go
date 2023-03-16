package controllers

import (
	"chat/internal/entities"

	"github.com/gofiber/fiber/v2"
)

type chatServiceHandler struct {
	ChatService entities.ChatService
}

func NewChatServiceHandler(app *fiber.App, c entities.ChatService) entities.ChatHandler {
	handler := &chatServiceHandler{ChatService: c}
	app.Post("/chatroom", handler.CreateRoomHandler)
	app.Post("/user/register", handler.RegisterHandler)
	app.Post("/user/login", handler.LoginHandler)
	return handler
}

func (chat *chatServiceHandler) CreateRoomHandler(c *fiber.Ctx) error {
	name := c.FormValue("room_name")
	id, err := chat.ChatService.CreateRoom(c.Context(), name)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"chatroom_id": id})
}

func (chat *chatServiceHandler) RegisterHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := chat.ChatService.Register(c.Context(), username, password)
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (chat *chatServiceHandler) LoginHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	access_token, refresh_token, err := chat.ChatService.Login(c.Context(), username, password)
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}
