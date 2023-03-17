package usecases

import (
	"chat/internal/entities"
	"context"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type chatService struct {
	repo     entities.ChatRepository
	rabbit   *amqp.Connection
	channels map[uuid.UUID]entities.ChatRoom
}

func NewChatService(rabbit_con *amqp.Connection, repo entities.ChatRepository) entities.ChatService {
	return &chatService{rabbit: rabbit_con, channels: make(map[uuid.UUID]entities.ChatRoom), repo: repo}
}

func (c *chatService) CreateRoom(ctx context.Context, name string) (*uuid.UUID, error) {
	ch, err := c.rabbit.Channel()
	if err != nil {
		return nil, err
	}
	id, err := c.repo.CreateRoom(ctx, name)
	if err != nil {
		return nil, err
	}
	c.channels[*id] = entities.ChatRoom{Name: name, Channel: ch}
	return id, nil
}
