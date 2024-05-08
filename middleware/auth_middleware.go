package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("secret-key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// Hilangkan respons jika tidak ada token
			if c.Request.URL.Path == "/users/register" {
				c.Next()
				return
			}
			if c.Request.URL.Path == "/users/login" {
				c.Next()
				return
			}
			if c.Request.URL.Path == "/users/login" {
				c.Next()
				return
			}
			c.JSON(401, gin.H{"error": "Authorization token not provided"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Menyimpan claims token JWT dalam context
		c.Set("claims", token.Claims)
		c.Next()
	}
}
