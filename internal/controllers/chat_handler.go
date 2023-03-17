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
