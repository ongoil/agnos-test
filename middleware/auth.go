package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "401",
				"message":    "authorization header is required",
				"message_th": "กรุณาระบุ Authorization Token",
			})
			c.Abort()
			return
		}

		// Check Bearer token
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "401",
				"message":    "invalid authorization format",
				"message_th": "รูปแบบ Authorization ไม่ถูกต้อง",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return signingKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "401",
				"message":    "invalid or expired token",
				"message_th": "Token ไม่ถูกต้องหรือหมดอายุ",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "message": "invalid token claims"})
			c.Abort()
			return
		}
		staffID, staffOK := claims["staff_id"].(string)
		hospitalID, hospitalOK := claims["hospital_id"].(string)
		if !staffOK || !hospitalOK || staffID == "" || hospitalID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "message": "invalid token claims"})
			c.Abort()
			return
		}
		c.Set("staff_id", staffID)
		c.Set("hospital_id", hospitalID)
		c.Set("token", token)

		c.Next()
	}
}
