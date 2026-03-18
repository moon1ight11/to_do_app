package usershandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
)

// получение данных пользователя
func (u *UserHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.logger.Error("Error in GetUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in GetUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in GetUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.logger.Info("User getted successfully", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		u.logger.Error("Error in UpdateUser ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// обновляем нужные поля
	err = u.userService.UpdateUser(ctx, updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in UpdateUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем обновленного пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.logger.Info("User updated successfully", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"updated_user": user})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.logger.Error("Error in DeleteUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// удаляем пользователя
	err = u.userService.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in DeleteUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in DeleteUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User deleted successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
