package utils

import (
	"net/http"
	"strings"
)

/*
GetAuthorizationHeader extracts the authorization token from the request's Authorization header.

Params:

	r *http.Request // The HTTP request.

Returns:

	typeToken string // The type of the authorization token (e.g., Bearer).
	token string // The actual authorization token.

Note: If the Authorization header is missing or invalid, an empty string is returned for both typeToken and token.

Example:

	typeToken, token := GetAuthorizationHeader(r)
	if typeToken == "Bearer" && token!= "" {
	  // Perform authentication and authorization checks
	} else {
	  http.Error(w, "Invalid or missing authorization token", http.StatusUnauthorized)
	}
*/
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
