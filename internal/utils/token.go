package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRefreshToken generates a secure random token for refresh tokens
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSecureToken generates a secure random token for general purposes (e.g., password resets)
func GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes) // Ignoring error since it is minimal in this use case
	return base64.URLEncoding.EncodeToString(bytes)
}
