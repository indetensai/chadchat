package entities

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

type UserRepository interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
	Validation(tokenstring string) (*jwt.Token, error)
}

type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
	Validation(tokenstring string) (*jwt.Token, error)
}

type UserHandler interface {
	RegisterHandler(c *fiber.Ctx) error
	LoginHandler(c *fiber.Ctx) error
}
