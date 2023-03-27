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
	Time   int64
	Limit  int64
	Offset int64
	RoomID uuid.UUID
}

type ChatRoom struct {
	Name   string
	RoomID uuid.UUID
}

type ChatHistory struct {
	Content   string
	CreatedAt time.Time
	Username  string
}

type ChatMessage struct {
	Username string    `json:"username"`
	Content  string    `json:"content"`
	SentAt   time.Time `json:"sent_at"`
	RoomID   uuid.UUID `json:"room_id"`
}

type ChatRepository interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	GetRoomByID(ctx context.Context, room_id uuid.UUID) (*ChatRoom, error)
	CreateMessage(ctx context.Context, message WriteMessageInput) error
	GetHistory(ctx context.Context, content GetHistoryInput) ([]ChatHistory, error)
	GetRooms(ctx context.Context) ([]ChatRoom, error)
}

type ChatService interface {
	CreateRoom(ctx context.Context, name string) (*uuid.UUID, error)
	CheckRoom(ctx context.Context, room_name uuid.UUID) error
	WriteMessage(ctx context.Context, message WriteMessageInput) error
	GetHistory(ctx context.Context, content GetHistoryInput) ([]ChatHistory, error)
	GetRooms(ctx context.Context) ([]ChatRoom, error)
}
