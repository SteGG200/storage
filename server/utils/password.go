package utils

import (
	"errors"

	"github.com/SteGG200/storage/logger"
	"golang.org/x/crypto/bcrypt"
)

const cost = 15

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil && !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		logger.ErrorLogger.Print(err)
	}

	return err == nil
}
