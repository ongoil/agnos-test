package middleware

import (
	"time"

	"testing"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)

	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func Add(a int, b int) int {
	return a + b
}

var jwtSecret = []byte("your-secret-key")

func GenerateToken(staffID uuid.UUID, hospitalID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"staff_id":    staffID.String(),
		"hospital_id": hospitalID.String(),
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
