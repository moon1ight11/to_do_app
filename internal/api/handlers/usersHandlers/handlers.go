package usershandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"todoapp/internal/api/models"
)

// получение данных пользователя
func (u *UserHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		u.logger.Error("Error in GetUser: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		u.logger.Error("Error in GetUser: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in GetUser:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in GetUser:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.logger.Info("User %v getted successfully", user.Id)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		u.logger.Error("Error in UpdateUser: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		u.logger.Error("Error in UpdateUser: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		u.logger.Error("Error in UpdateUser ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// обновляем нужные поля
	err := u.userService.UpdateUser(ctx, updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in UpdateUser:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in UpdateUser:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем обновленного пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.logger.Info("User %v updated successfully", user.Id)

	c.JSON(http.StatusOK, gin.H{"updated_user": user})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		u.logger.Error("Error in DeleteUser: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		u.logger.Error("Error in DeleteUser: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// удаляем пользователя
	err := u.userService.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in DeleteUser: %w", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in DeleteUser: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User %v deleted successfully", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
