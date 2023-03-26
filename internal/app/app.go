package app

import (
	"chat/internal/controllers"
	"chat/internal/controllers/auth"
	"chat/internal/repository"
	"chat/internal/usecases"
	"context"
	"crypto/rsa"
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
)

func get_private_key(filename string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(filename)
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
	return privatekey
}

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
	/*client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping(); err != nil {
		log.Fatal("redis kaput")
	}*/
	app := fiber.New()

	user_repository := repository.NewUserRepository(
		pgx_con,
		get_private_key("access_private.pem"),
		get_private_key("refresh_private.pem"),
	)
	user_service := usecases.NewUserService(user_repository)
	controllers.NewUserServiceHandler(app, user_service)

	app.Use(auth.New(user_repository))

	chat_repository := repository.NewChatRepository(pgx_con)
	chat_service := usecases.NewChatService(chat_repository)
	controllers.NewChatServiceHandler(app, chat_service)

	app.Listen(":8080")
}
