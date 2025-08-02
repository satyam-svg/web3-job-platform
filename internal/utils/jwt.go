package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID uuid.UUID) string {
	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(), // Convert UUID to string
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}
	return ss
}
