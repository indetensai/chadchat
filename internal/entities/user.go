package entities

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

type TokenCredentials struct {
	UserID   string
	Username string
}

type UserRepository interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
	Validation(tokenstring string) (*TokenCredentials, error)
	Refresh(tokenstring string) (*string, *string, error)
}

type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*string, *string, error)
	Validation(tokenstring string) (*TokenCredentials, error)
	Refresh(tokenstring string) (*string, *string, error)
}
