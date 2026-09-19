package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte(func() string {
	if value := os.Getenv("JWT_SECRET"); value != "" {
		return value
	}
	return "development-only-change-me"
}())

func GenerateToken(staffID uuid.UUID, hospitalID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"staff_id": staffID.String(), "hospital_id": hospitalID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(), "iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}
