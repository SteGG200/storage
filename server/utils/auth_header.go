package utils

import (
	"net/http"
	"strings"
)

func GetAuthorizationHeader(r *http.Request) (typeToken string, token string) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", ""
	}

	slice := strings.SplitN(authHeader, " ", 2)

	if len(slice) < 2 {
		return "", ""
	}

	typeToken = slice[0]
	token = slice[1]

	return
}
