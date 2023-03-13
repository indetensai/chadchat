package app

import (
	"chat/internal/controllers"
	"chat/internal/usecases"
	"fmt"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Run() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Print("baza kaput")
	}
	app := fiber.New()
	defer conn.Close()
	chat_service := usecases.NewChatService(conn)
	controllers.NewChatServiceHandler(app, chat_service)
	app.Listen(":8080")
}
