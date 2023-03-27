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

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}

type UserService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (*Tokens, error)
	Validation(tokenstring string) (*TokenCredentials, error)
	Refresh(tokenstring string) (*Tokens, error)
}
