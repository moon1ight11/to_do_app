package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todoapp/internal/api/jwt"
	"todoapp/pkg/logger"
)

// мидлвар аутентификации
func Auth(jwtService jwt.TokenService, logger logger.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			logger.Error("Error in get value from cookie in AuthMiddle:", "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			logger.Error("Error in parse token in AuthMiddle:", "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Error("Error in AuthMiddle: Token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Set("UserId", *claims.UserId)
	}
}
