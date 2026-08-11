package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todoapp/internal/api/jwt"
	"todoapp/pkg/logger"
)

func Auth(jwtService jwt.TokenService, l logger.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("token")
		if err != nil {
			l.Error("middleware.Auth: get cookie", "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			l.Error("middleware.Auth: parse token", "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		if !token.Valid {
			l.Error("middleware.Auth: token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Set("UserId", *claims.UserId)
	}
}
