package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		userID, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)

		if strings.TrimSpace(userID) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)

		c.Next()
	}
}

func ExtractUser(c *gin.Context) (string, string, bool) {
	uid, ok1 := c.Get("user_id")
	email, ok2 := c.Get("email")
	if !ok1 || !ok2 {
		return "", "", false
	}
	return uid.(string), email.(string), true
}
