package utils

import "fmt"

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
			err := fmt.Errorf("Cannot convert path to string")
			return false, err
		}

		if currentPathStr == path {
			exist = true
			break
		}
	}

	return exist, nil
}
