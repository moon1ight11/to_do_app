package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"todoapp/internal/api/jwt"
)

// мидлвар аутентификации
func Auth(jwtService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			log.Println("Error in get value from cookie")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			log.Println("Error in parse token", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Set("UserId", claims.UserId)
	}
}
