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
	app.Post("/chatroom", handler.CreateChatRoomHandler)
	return handler
}

func (chat *chatServiceHandler) CreateChatRoomHandler(c *fiber.Ctx) error {
	name := c.FormValue("room_name")
	id, err := chat.ChatService.CreateChatRoom(name)
	if err != nil {
		return fiber.ErrInternalServerError //error handling must be here
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"chatroom_id": id})
}
