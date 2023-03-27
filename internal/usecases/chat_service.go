package usecases

import (
	"chat/internal/entities"
	"context"

	"github.com/google/uuid"
)

type chatService struct {
	repo entities.ChatRepository
}

func NewChatService(repo entities.ChatRepository) entities.ChatService {
	return &chatService{
		repo: repo,
	}
}

func (c *chatService) CreateRoom(ctx context.Context, name string) (*uuid.UUID, error) {
	id, err := c.repo.CreateRoom(ctx, name)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (c *chatService) WriteMessage(ctx context.Context, message entities.WriteMessageInput) error {
	return c.repo.WriteMessage(ctx, message)
}

func (c *chatService) CheckRoom(ctx context.Context, room_id uuid.UUID) error {
	return c.repo.CheckRoom(ctx, room_id)
}

func (c *chatService) GetHistory(
	ctx context.Context,
	content entities.GetHistoryInput,
) (*[]entities.ChatHistory, error) {
	return c.repo.GetHistory(ctx, content)
}

func (c *chatService) GetRooms(ctx context.Context) (*[]entities.ChatRoom, error) {
	return c.repo.GetRooms(ctx)
}
