package utils

import (
	"fmt"
	"os"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
	"github.com/golang-jwt/jwt/v5"
)

/*
GenerateToken generates a JSON Web Token (JWT) with the provided data.

Params:

	data map[string]any // The data to be included in the JWT claims.

Returns:

	typeToken string // The generated JWT token.
	err error // An error if any occurred during the token generation process.
*/
func GenerateToken(data map[string]any) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		logger.ErrorLogger.Print(exception.ErrNotFoundJWTSecret)
		return "", exception.ErrNotFoundJWTSecret
	}

	tokenString, err = token.SignedString([]byte(jwtSecret))

	if err != nil {
		logger.ErrorLogger.Print(err)
		return "", err
	}

	return tokenString, nil
}

/*
ParseToken parses a JSON Web Token (JWT) and returns its data.

Params:

	tokenString string // The JWT token to parse.

Returns:

	map[string]any // The parsed data.
	err error // An error if any occurred during the token parsing process.
*/
func ParseToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			logger.ErrorLogger.Print(err)
			return nil, err
		}

		jwtSecret, ok := os.LookupEnv("JWT_SECRET")

		if !ok {
			logger.ErrorLogger.Print(exception.ErrNotFoundJWTSecret)
			return nil, exception.ErrNotFoundJWTSecret
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		err := fmt.Errorf("unexpected claims type")
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	return claims, nil
}
