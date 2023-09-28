package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
