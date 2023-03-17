package usecases

import (
	"chat/internal/entities"
	"context"
)

type userService struct {
	repo entities.UserRepository
}

func NewUserService(repo entities.UserRepository) entities.UserService {
	return &userService{repo: repo}
}
func (u *userService) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	return u.repo.Register(ctx, username, password)
}

func (u *userService) Login(
	ctx context.Context,
	username string,
	password string,
) (*string, *string, error) {
	access_token, refresh_token, err := u.repo.Login(ctx, username, password)
	if err != nil {
		return nil, nil, err
	}
	return access_token, refresh_token, nil
}

func (u *userService) Validation(tokenstring string) (*string, *string, error) {
	return u.repo.Validation(tokenstring)
}

func (u *userService) Refresh(tokenstring string) (*string, *string, error) {
	return u.repo.Refresh(tokenstring)
}
