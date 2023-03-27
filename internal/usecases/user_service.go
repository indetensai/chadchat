package usecases

import (
	"chat/internal/entities"
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo        entities.UserRepository
	access_key  *rsa.PrivateKey
	refresh_key *rsa.PrivateKey
}

type Claims struct {
	jwt.StandardClaims
	UserID   string
	Username string
}

func NewUserService(
	repo entities.UserRepository,
	access_key *rsa.PrivateKey,
	refresh_key *rsa.PrivateKey,
) entities.UserService {
	return &userService{repo: repo, access_key: access_key, refresh_key: refresh_key}
}

func (u *userService) GenerateJWTAccessToken(user_id string, username string) (*string, error) {
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

func (u *userService) GenerateJWTRefreshToken(user_id string, username string) (*string, error) {
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

func (u *userService) Validation(tokenstring string) (*entities.TokenCredentials, error) {
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
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, entities.ErrInvalidCredentials
		}
		return &entities.TokenCredentials{UserID: claims.UserID, Username: claims.Username}, nil
	} else {
		return nil, err
	}
}

func (u *userService) Refresh(tokenstring string) (*entities.Tokens, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, check := t.Method.(*jwt.SigningMethodRSA); !check {
			return nil, entities.ErrInvalidCredentials
		}
		return &u.refresh_key.PublicKey, nil
	})
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}
	if token.Valid {
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, entities.ErrInvalidCredentials
		}
		access_token, err := u.GenerateJWTAccessToken(claims.UserID, claims.Username)
		if err != nil {
			return nil, err
		}
		refresh_token, err := u.GenerateJWTRefreshToken(claims.UserID, claims.Username)
		if err != nil {
			return nil, err
		}
		return &entities.Tokens{AccessToken: *access_token, RefreshToken: *refresh_token}, nil
	} else {
		return nil, err
	}
}

func (u *userService) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	return u.repo.CreateUser(ctx, username, password)
}

func (u *userService) Login(
	ctx context.Context,
	username string,
	password string,
) (*entities.Tokens, error) {

	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}
	access_token, err := u.GenerateJWTAccessToken(user.ID.String(), username)
	if err != nil {
		return nil, err
	}
	refresh_token, err := u.GenerateJWTRefreshToken(user.ID.String(), username)

	if err != nil {
		return nil, err
	}
	return &entities.Tokens{RefreshToken: *refresh_token, AccessToken: *access_token}, nil
}
