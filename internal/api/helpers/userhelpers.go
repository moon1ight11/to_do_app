package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получение id пользователя из контекста
func GetUserIdFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exist := c.Get("UserId")
	if !exist {
		return uuid.Nil, fmt.Errorf("error in GetUserIdFromContext: user ID not found in context")
	}

	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("error in GetUserIdFromContext: invalid user ID type")
	}

	return userId, nil
}
