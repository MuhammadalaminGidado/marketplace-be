package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateSessionToken() (string, string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}

	token := hex.EncodeToString(b)

	hash := sha256.Sum256([]byte(token))

	return token, hex.EncodeToString(hash[:]), nil
}

func HashSessionToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate csrf token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
