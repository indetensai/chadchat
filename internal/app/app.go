package app

import (
	"chat/internal/controllers"
	"chat/internal/usecases"
	"chat/internal/usecases/repository"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Run() {
	godotenv.Load(".env")
	m, err := migrate.New("file://../migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		fmt.Print(err)
	}
	pgx_con, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("baza kaput")
	}
	if err = pgx_con.Ping(context.Background()); err != nil {
		log.Fatal("baza kaput")
	}
	rabbit_con, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	if err != nil {
		log.Fatal("rabbit kaput")
	}
	defer rabbit_con.Close()
	/*client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping(); err != nil {
		log.Fatal("redis kaput")
	}*/
	privateKeyFile, err := os.Open("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyPEM, _ := pem.Decode(privateKeyBytes)
	privatekey, err := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	chat_repository := repository.NewChatRepository(pgx_con)
	chat_service := usecases.NewChatService(rabbit_con, chat_repository)
	controllers.NewChatServiceHandler(app, chat_service)

	user_repository := repository.NewUserRepository(pgx_con, privatekey)
	user_service := usecases.NewUserService(user_repository)
	controllers.NewUserServiceHandler(app, user_service)

	app.Listen(":8080")
}
