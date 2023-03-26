package controllers

import (
	"chat/internal/entities"
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type chatUser struct {
	isClosing bool
	mu        sync.Mutex
}

type roomConnection struct {
	Connection *websocket.Conn
	RoomID     uuid.UUID
}

type chatServiceHandler struct {
	ChatService entities.ChatService
	listeners   map[uuid.UUID]map[*websocket.Conn]*chatUser
	register    chan roomConnection
	broadcast   chan entities.ChatMessage
	unregister  chan roomConnection
}

func NewChatServiceHandler(app *fiber.App, c entities.ChatService) {
	handler := &chatServiceHandler{
		ChatService: c,
		listeners:   make(map[uuid.UUID]map[*websocket.Conn]*chatUser),
		register:    make(chan roomConnection),
		broadcast:   make(chan entities.ChatMessage),
		unregister:  make(chan roomConnection),
	}
	app.Post("/chatroom", handler.CreateRoomHandler)
	go handler.RunHub()
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/:room_id<guid>", websocket.New(handler.GetWebsocketConnection))
	app.Get("/chat/:room_id<guid>/history", handler.GetHistory)
}

func (chat *chatServiceHandler) RunHub() {
	for {
		select {
		case connection := <-chat.register:
			chat.listeners[connection.RoomID][connection.Connection] = &chatUser{}
			log.Println("connection registered")

		case message := <-chat.broadcast:
			log.Println("message received:", message)
			for connection, c := range chat.listeners[message.RoomID] {
				go func(connection *websocket.Conn, c *chatUser) {
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}
					if err := connection.WriteMessage(
						websocket.TextMessage,
						[]byte(message.Username+": "+message.Content),
					); err != nil {
						c.isClosing = true
						log.Println("write error:", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						chat.unregister <- roomConnection{
							RoomID:     message.RoomID,
							Connection: connection,
						}
					}
				}(connection, c)
			}

		case connection := <-chat.unregister:
			delete(chat.listeners[connection.RoomID], connection.Connection)

			log.Println("connection unregistered")
		}
	}
}

func (chat *chatServiceHandler) CreateRoomHandler(c *fiber.Ctx) error {
	name := c.FormValue("room_name")
	id, err := chat.ChatService.CreateRoom(c.Context(), name)
	if err != nil {
		return error_handling(c, err)
	}
	chat.listeners[*id] = make(map[*websocket.Conn]*chatUser)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"chatroom_id": id})
}

func (chat *chatServiceHandler) GetWebsocketConnection(c *websocket.Conn) {
	room_id, err := uuid.Parse(c.Params("room_id"))
	if err != nil {
		log.Fatal(err)
	}
	user_id, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		log.Fatal(err)
	}
	_, ok := chat.listeners[room_id]
	if !ok {
		err = chat.ChatService.CheckRoom(context.Background(), room_id)
		if err != nil {
			log.Fatal(err)
		}
		chat.listeners[room_id] = make(map[*websocket.Conn]*chatUser)
	}
	username := c.Locals("username").(string)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		chat.unregister <- roomConnection{RoomID: room_id, Connection: c}
		c.Close()
	}()

	chat.register <- roomConnection{RoomID: room_id, Connection: c}
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println(err)
			}

			return
		}

		if messageType == websocket.TextMessage {
			chat.broadcast <- entities.ChatMessage{
				Content:  string(message),
				Username: username,
				SentAt:   time.Now().Format("15:04:05"),
				RoomID:   room_id,
			}
			err = chat.ChatService.WriteMessage(context.Background(), entities.WriteMessageInput{
				Content:   string(message),
				UserID:    user_id,
				CreatedAt: time.Now(),
				RoomID:    room_id,
				Username:  username,
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}

func (chat *chatServiceHandler) GetHistory(c *fiber.Ctx) error {
	room_id, err := uuid.Parse(c.Params("room_id"))
	time, _ := strconv.ParseInt(c.Query("time"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		return err
	}
	history, err := chat.ChatService.GetHistory(c.Context(), entities.GetHistoryInput{
		RoomID: room_id,
		Time:   time,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"history": history})
}
