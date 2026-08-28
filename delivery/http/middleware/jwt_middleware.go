package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึงค่าจาก Header: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// 2. ตรวจสอบความถูกต้องของ Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 3. แกะข้อมูล (Claims) ออกมา
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// ฝากข้อมูลที่จำเป็นไว้ใน Gin Context เพื่อให้ Handler ตัวถัดไปหยิบไปใช้
			c.Set("user_id", claims["user_id"])
			c.Set("hospital_id", claims["hospital_id"])
			c.Next() // อนุญาตให้ไปทำงานที่ Handler ถัดไปได้
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
	}
}