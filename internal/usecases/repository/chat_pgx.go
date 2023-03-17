package repository

import (
	"chat/internal/entities"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type chatRepository struct {
	db *pgx.Conn
}

func NewChatRepository(db *pgx.Conn) entities.ChatRepository {
	return &chatRepository{db: db}
}

func (c *chatRepository) CreateRoom(ctx context.Context, name string) (*uuid.UUID, error) {
	var id uuid.UUID
	err := c.db.QueryRow(
		ctx,
		"INSERT INTO rooms (room_name) VALUES ($1) RETURNING room_id",
		name,
	).Scan(&id)
	if err != nil {
		return nil, entities.ErrDuplicate
	}
	return &id, nil
}
