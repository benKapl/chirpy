package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header not found")
	}

	bearerPrefix := "Bearer "

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("Invalide authorization header format")
	}

	if authHeader == bearerPrefix {
		return "", fmt.Errorf("No token provided")
	}

	bearerToken := strings.Replace(authHeader, bearerPrefix, "", 1)
	return bearerToken, nil
}
