package repository

import (
	"chat/internal/entities"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return nil, entities.ErrDuplicate
		}
		return nil, err
	}
	return &id, nil
}

func (c *chatRepository) WriteMessage(
	ctx context.Context,
	message entities.WriteMessageInput,
) error {
	_, err := c.db.Exec(
		ctx,
		"INSERT INTO messages (user_id,room_id,content,sent_at,username) VALUES ($1,$2,$3,$4,$5)",
		message.UserID,
		message.RoomID,
		message.Content,
		message.CreatedAt,
		message.Username,
	)
	return err
}

func (c *chatRepository) CheckRoom(ctx context.Context, room_id uuid.UUID) error {
	_, err := c.db.Exec(
		ctx,
		"SELECT room_name FROM rooms WHERE room_id=$1",
		room_id,
	)
	if err != nil {
		return entities.ErrNotFound
	}
	return nil
}

func (c *chatRepository) GetHistory(
	ctx context.Context,
	content entities.GetHistoryInput,
) (*[]entities.ChatHistory, error) {
	rows, err := c.db.Query(
		ctx,
		"SELECT (content,sent_at,username) FROM messages WHERE sent_at<=$1 AND room_id=$2 ORDER BY sent_at DESC LIMIT $3 OFFSET $4",
		time.Unix(content.Time, 0),
		content.RoomID,
		content.Limit,
		content.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var history []entities.ChatHistory
	for rows.Next() {
		var r entities.ChatHistory
		err = rows.Scan(&r)
		if err != nil {
			return nil, err
		}
		history = append(history, r)
	}
	return &history, nil
}

func (c *chatRepository) GetRooms(ctx context.Context) (*[]entities.ChatRoom, error) {
	rows, err := c.db.Query(ctx, "SELECT (room_name,room_id) FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rooms []entities.ChatRoom
	for rows.Next() {
		var r entities.ChatRoom
		err = rows.Scan(&r)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return &rooms, nil
}
