package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func CheckPasswordHash(hash, password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}
func CompareWithDummyHash(password string) {
	_ = bcrypt.CompareHashAndPassword(
		[]byte("$2a$10$N9qo8uLOickgx2ZMRZo5i.ej6tqT0h1Qyq1Z5l6b5o5e5e5e5e5e"),
		[]byte(password),
	)
}
