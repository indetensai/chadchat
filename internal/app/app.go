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
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func getPrivateKey(filename string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open private key file: `", filename, "` ", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := ioutil.ReadAll(privateKeyFile)
	if err != nil {
		log.Fatal("failed to read private key file: `", filename, "` ", err)
	}

	privateKeyPEM, _ := pem.Decode(privateKeyBytes)
	privatekey, err := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		log.Fatal("failed to parse private key file: `", filename, "` ", err)
	}
	return privatekey
}

func Run(databaseURL string, port string) {
	pgx_con, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	if err = pgx_con.Ping(context.Background()); err != nil {
		log.Print("warn - database isn't respodning to ping: ", err)
	}
	app := fiber.New()

	user_repository := repository.NewUserRepository(
		pgx_con,
		getPrivateKey("access_private.pem"),
		getPrivateKey("refresh_private.pem"),
	)
	user_service := usecases.NewUserService(user_repository)
	controllers.NewUserServiceHandler(app, user_service)

	app.Use(skip.New(auth.New(user_repository), func(c *fiber.Ctx) bool {
		return c.Path() == "/chat/rooms"
	}))

	chat_repository := repository.NewChatRepository(pgx_con)
	chat_service := usecases.NewChatService(chat_repository)
	controllers.NewChatServiceHandler(app, chat_service)

	app.Listen(":" + port)
}
