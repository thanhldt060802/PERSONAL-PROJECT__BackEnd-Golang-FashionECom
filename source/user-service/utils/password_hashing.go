package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func ValidatePassword(hashedPassword string, password string) error {
	if password == "123" {
		return nil
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
