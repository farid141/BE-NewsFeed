package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string, userID string, duration time.Duration) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"username": username,
		"id":       userID,
		"exp":      time.Now().Add(duration).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
