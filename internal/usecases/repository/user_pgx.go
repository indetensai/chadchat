package repository

import (
	"chat/internal/entities"
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db          *pgx.Conn
	access_key  *rsa.PrivateKey
	refresh_key *rsa.PrivateKey
}

type Claims struct {
	jwt.StandardClaims
	UserID   string
	Username string
}

func NewUserRepository(
	db *pgx.Conn,
	access_key *rsa.PrivateKey,
	refresh_key *rsa.PrivateKey,
) entities.UserRepository {
	return &userRepository{db: db, access_key: access_key, refresh_key: refresh_key}
}

func (u *userRepository) GenerateJWTAccessToken(user_id string, username string) (*string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		UserID:   user_id,
		Username: username,
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := tokenString.SignedString(u.access_key)
	if err != nil {
		return nil, err
	}
	return &token, err
}

func (u *userRepository) GenerateJWTRefreshToken(user_id string, username string) (*string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 28).Unix(),
		},
		UserID:   user_id,
		Username: username,
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := tokenString.SignedString(u.refresh_key)
	if err != nil {
		return nil, err
	}
	return &token, err
}
func (u *userRepository) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	tx, err := u.db.Begin(ctx)
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

func (u *userRepository) Login(
	ctx context.Context,
	username string,
	password string,
) (*string, *string, error) {
	tx, err := u.db.Begin(ctx)
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
	access_token, err := u.GenerateJWTAccessToken(result.ID.String(), result.Username)
	if err != nil {
		return nil, nil, err
	}
	refresh_token, err := u.GenerateJWTRefreshToken(result.ID.String(), result.Username)

	if err != nil {
		return nil, nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, nil, err
	}
	return access_token, refresh_token, nil
}

func (u *userRepository) Validation(tokenstring string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, check := t.Method.(*jwt.SigningMethodRSA); !check {
			return nil, entities.ErrInvalidCredentials
		}
		return &u.access_key.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}

func (u *userRepository) Refresh(tokenstring string) (*string, *string, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, check := t.Method.(*jwt.SigningMethodRSA); !check {
			return nil, entities.ErrInvalidCredentials
		}
		return &u.refresh_key.PublicKey, nil
	})
	if err != nil {
		return nil, nil, entities.ErrInvalidCredentials
	}
	if token.Valid {
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, nil, entities.ErrInvalidCredentials
		}
		access_token, err := u.GenerateJWTAccessToken(claims.UserID, claims.Username)
		if err != nil {
			return nil, nil, err
		}
		refresh_token, err := u.GenerateJWTRefreshToken(claims.UserID, claims.Username)
		if err != nil {
			return nil, nil, err
		}
		return access_token, refresh_token, nil
	} else {
		return nil, nil, err
	}
}
