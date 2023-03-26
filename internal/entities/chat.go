package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type WriteMessageInput struct {
	RoomID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	Username  string
}

type GetHistoryInput struct {
	Time   int64     `json:"time"`
	Limit  int64     `json:"limit"`
	Offset int64     `json:"offset"`
	RoomID uuid.UUID `json:"room_id"`
}

type GetHistoryOutput struct {
	Content   string
	CreatedAt time.Time
	Username  string
}

type ChatMessage struct {
	Username string
	Content  string
	SentAt   string
	RoomID   uuid.UUID
}

type ChatRepository interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	CheckRoom(ctx context.Context, room_name uuid.UUID) error
	WriteMessage(ctx context.Context, message WriteMessageInput) error
	GetHistory(ctx context.Context, content GetHistoryInput) (*[]GetHistoryOutput, error)
}

type ChatService interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	CheckRoom(ctx context.Context, room_name uuid.UUID) error
	WriteMessage(ctx context.Context, message WriteMessageInput) error
	GetHistory(ctx context.Context, content GetHistoryInput) (*[]GetHistoryOutput, error)
}
