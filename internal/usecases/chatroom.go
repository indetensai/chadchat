package usecases

import (
	"chat/internal/entities"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type chatService struct {
	rabbit   *amqp.Connection
	channels map[uuid.UUID]entities.ChatRoom
}

func NewChatService(conn *amqp.Connection) entities.ChatService {
	return &chatService{rabbit: conn, channels: make(map[uuid.UUID]entities.ChatRoom)}
}

func (c *chatService) CreateChatRoom(name string) (*uuid.UUID, error) {
	ch, err := c.rabbit.Channel()
	if err != nil {
		return nil, err
	}
	id := uuid.New()
	c.channels[id] = entities.ChatRoom{Name: name, Channel: ch}
	return &id, nil
}
