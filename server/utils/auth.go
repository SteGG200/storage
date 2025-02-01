package utils

import "slices"

func VerifyAuthorizationToken(tokenString string, path string) (bool, error) {
	data, err := ParseToken(tokenString)

	if err != nil {
		return false, err
	}

	exist := slices.Contains(data["path"].([]string), path)
	return exist, nil
}
