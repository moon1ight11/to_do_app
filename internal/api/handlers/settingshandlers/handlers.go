package settingshandlers

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

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		log.Println("Error in getSettings: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		log.Println("Error in getSettings: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// находим настройки по id
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Error in getSettings:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in getSettings:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		log.Println("Error in updateSettings: user ID not found in context")
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		log.Println("Error in updateSettings: invalid user ID type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedSettings models.UpdatedSettings

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		log.Println("Error in updateSettings ShouldBindJSON: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// изменяем настройки
	err := s.settingsService.UpdateSettings(ctx, userId, updatedSettings.TimeDuration, updatedSettings.UserTz)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Error in updateSettings: %w", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in updateSettings: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем измененные настройки
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		log.Println("Error in updateSettings: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated_settings": settings})
}
