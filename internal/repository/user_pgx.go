package repository

import (
	"chat/internal/entities"
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) entities.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(
	ctx context.Context,
	username string,
	password string,
) error {
	password_raw := []byte(password)
	password_hashed, err := bcrypt.GenerateFromPassword(password_raw, 10)
	if err != nil {
		return err
	}
	_, err = u.db.Exec(
		ctx,
		"INSERT INTO users (username,password) VALUES ($1,$2)",
		username,
		string(password_hashed),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return entities.ErrDuplicate
		}
		return err
	}
	return nil
}

func (u *userRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*entities.User, error) {
	var result entities.User
	err := u.db.QueryRow(
		ctx,
		"SELECT * FROM users WHERE username=$1",
		username,
	).Scan(&result.ID, &result.Username, &result.Password)
	if err != nil {
		return nil, entities.ErrNotFound
	}
	return &result, nil
}
