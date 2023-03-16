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

func (c *chatService) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	return c.repo.Register(ctx, username, password)
}

func (c *chatService) Login(
	ctx context.Context,
	username string,
	password string,
) (*string, *string, error) {
	access_token, refresh_token, err := c.repo.Login(ctx, username, password)
	if err != nil {
		return nil, nil, err
	}
	return access_token, refresh_token, nil
}
