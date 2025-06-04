package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword)); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	return nil
}

func HashToken(token, hashSecret string) (hashedToken string, err error) {
	h := hmac.New(sha256.New, []byte(hashSecret))

	_, err = h.Write([]byte(token))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func VerifyToken(hashedToken, token, hashSecret string) bool {
	expectedHash, err := HashToken(token, hashSecret)
	if err != nil {
		log.Println(err)
		return false
	}

	return hmac.Equal([]byte(hashedToken), []byte(expectedHash))
}
