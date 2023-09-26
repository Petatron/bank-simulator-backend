package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hash a password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("there is an error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword check if a password is correct
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
