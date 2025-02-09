package utils

import "fmt"

/*
VerifyAuthorizationToken checks if the provided authorization token is valid and has access to the current path.

Params:

	tokenString string // The authorization token to verify.
	path string // The path to check access for.
*/
func VerifyAuthorizationToken(tokenString string, path string) (bool, error) {
	if len(tokenString) == 0 {
		return false, nil
	}
	data, err := ParseToken(tokenString)

	if err != nil {
		return false, err
	}

	exist := false

	for _, currentPath := range data["path"].([]any) {
		currentPathStr, ok := currentPath.(string)

		if !ok {
			err := fmt.Errorf("cannot convert path to string")
			return false, err
		}

		if currentPathStr == path {
			exist = true
			break
		}
	}

	return exist, nil
}
