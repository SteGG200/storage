package exception

import "errors"

var (
	ErrNotADirectory = errors.New("not a directory")
)
