package main

import (
	"chat/internal/app"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	databaseURL := os.Getenv("DATABASE_URL")
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatal("failed to read migration ", err)
	}
	if err := m.Up(); err != nil {
		log.Print("warn - migration failed: ", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(databaseURL, port)
}
