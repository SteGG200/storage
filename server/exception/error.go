package exception

import "errors"

var (
	ErrNotADirectory     = errors.New("not a directory")
	ErrNotAFile          = errors.New("not a file")
	ErrNotFoundJWTSecret = errors.New("not found env variable JWT_SECRET")
	ErrInvalidToken      = errors.New("invalid token")
)
