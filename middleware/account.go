package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func signingKey() []byte {
	if value := os.Getenv("JWT_SECRET"); value != "" {
		return []byte(value)
	}
	return []byte("development-only-change-me")
}

func tokenLifetime() time.Duration {
	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		if hours, err := time.ParseDuration(value + "h"); err == nil && hours > 0 {
			return hours
		}
	}
	return 24 * time.Hour
}

func GenerateToken(staffID uuid.UUID, hospitalID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"staff_id": staffID.String(), "hospital_id": hospitalID.String(),
		"exp": time.Now().Add(tokenLifetime()).Unix(), "iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey())
}
