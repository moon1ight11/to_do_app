package usershandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"todoapp/internal/api/models"
)

// получение данных пользователя
func (u *UserHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		log.Println("Error in get user: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		log.Println("Error in get user: invalid user ID type")
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
			log.Println("Error in get user: %w", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in get user: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		log.Println("Error in update user: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		log.Println("Error in update user: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Println("Error in update user ShouldBindJSON: %w", err)
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
			log.Println("Error in update user: %w", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in update user: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем обновленного пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		log.Println("Error in update user: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated_user": user})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		log.Println("Error in dalete user: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		log.Println("Error in delete user: invalid user ID type")
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
			log.Println("Error in delete user: %w", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in delete user: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
