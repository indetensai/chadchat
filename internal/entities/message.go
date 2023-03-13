package entities

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	ProducerID uuid.UUID
	ID         uuid.UUID
	Content    string
}

type ChatRoom struct {
	Channel *amqp.Channel
	Name    string
}

type ChatService interface {
	CreateChatRoom(name string) (*uuid.UUID, error)
}

type ChatHandler interface {
	CreateChatRoomHandler(c *fiber.Ctx) error
}
