// internal/domain/error.go

package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidAuthToken  = errors.New("invalid auth token")
)
