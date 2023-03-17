package entities

import "errors"

var (
	ErrEmptySession       = errors.New("empty session")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
	ErrNotAuthorized      = errors.New("don't have permission")
	ErrDuplicate          = errors.New("aldready exists")
)
