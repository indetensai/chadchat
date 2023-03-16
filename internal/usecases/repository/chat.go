package repository

import (
	"chat/internal/entities"
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ChatRepository struct {
	db  *pgx.Conn
	key *rsa.PrivateKey
}

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string
}

func NewChatRepository(db *pgx.Conn, key *rsa.PrivateKey) entities.ChatRepository {
	return &ChatRepository{db: db, key: key}
}

func (c *ChatRepository) GenerateJWTAccessToken(user_id string, username string) (*string, error) {
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(36000)).Unix(),
		},
		CustomClaims: map[string]string{
			"user_id":  user_id,
			"username": username,
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := tokenString.SignedString(c.key)
	if err != nil {
		return nil, err
	}
	return &token, err
}

func (c *ChatRepository) GenerateJWTRefreshToken(user_id string, username string) (*string, error) {
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(2419200)).Unix(),
		},
		CustomClaims: map[string]string{
			"user_id":  user_id,
			"username": username,
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := tokenString.SignedString(c.key)
	if err != nil {
		return nil, err
	}
	return &token, err
}

func (c *ChatRepository) CreateRoom(ctx context.Context, name string) (*uuid.UUID, error) {
	var id uuid.UUID
	err := c.db.QueryRow(
		ctx,
		"INSERT INTO rooms (room_name) VALUES ($1) RETURNING room_id",
		name,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (c *ChatRepository) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	password_raw := []byte(password)
	password_hashed, err := bcrypt.GenerateFromPassword(password_raw, 10)
	if err != nil {
		return err
	}
	var exists bool
	err = tx.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)",
		username,
	).Scan(&exists)
	if err != nil {
		return entities.ErrNotFound
	}
	if exists || (username == "" || password == "") {
		return entities.ErrDuplicate

	} else {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO users (username,password) VALUES ($1,$2)",
			username,
			string(password_hashed),
		)
		if err != nil {
			return entities.ErrDuplicate
		}
		tx.Commit(ctx)
		return nil
	}
}

func (c *ChatRepository) Login(
	ctx context.Context,
	username string,
	password string,
) (*string, *string, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)
	var result entities.User
	err = tx.QueryRow(
		ctx,
		"SELECT * FROM users WHERE username=$1",
		username,
	).Scan(&result.ID, &result.Username, &result.Password)
	if err != nil {
		return nil, nil, entities.ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return nil, nil, entities.ErrInvalidCredentials
	}
	access_token, err := c.GenerateJWTAccessToken(result.ID.String(), result.Username)
	if err != nil {
		return nil, nil, err
	}
	refresh_token, err := c.GenerateJWTRefreshToken(result.ID.String(), result.Username)

	if err != nil {
		return nil, nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, nil, err
	}
	return access_token, refresh_token, nil
}
