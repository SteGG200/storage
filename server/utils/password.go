package utils

import (
	"errors"

	"github.com/SteGG200/storage/logger"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

/*
HashPassword hashes the provided password using bcrypt.

Params:

	password string // The password to hash.
*/
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

/*
VerifyPassword verifies that the provided password matches the hashed password.

Params:

	password string // The provided password.
	hashedPassword string // The hashed password.
*/
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil && !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		logger.ErrorLogger.Print(err)
	}

	return err == nil
}
