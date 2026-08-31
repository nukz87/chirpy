package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKeyHeader := headers.Get("Authorization")
	if apiKeyHeader == "" {
		return "", errors.New("No auth header")
	}

	apiKey := strings.Split(apiKeyHeader, " ")
	if len(apiKey) < 2 || apiKey[0] != "ApiKey" {
		return "", errors.New("Malformed auth header")
	}

	return apiKey[1], nil
}
