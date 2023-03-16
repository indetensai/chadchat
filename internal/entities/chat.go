package entities

import (
	"context"

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

type ChatRepository interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
}

type ChatService interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
}

type ChatHandler interface {
	CreateRoomHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	LoginHandler(c *fiber.Ctx) error
}
