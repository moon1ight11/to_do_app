package settingshandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
)

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.logger.Error("Error in GetSettings:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// находим настройки по id
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.logger.Error("Error in GetSettings:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		s.logger.Error("Error in GetSettings:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Setting getted successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.logger.Error("Error in UpdateSettings:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var updatedSettings models.UpdatedSettings

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		s.logger.Error("Error in UpdateSettings ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// изменяем настройки
	err = s.settingsService.UpdateSettings(ctx, userId, updatedSettings.TimeDuration, updatedSettings.UserTz)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.logger.Error("Error in UpdateSettings:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		s.logger.Error("Error in UpdateSettings:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем измененные настройки
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		s.logger.Error("Error in UpdateSettings:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Settings updated successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_settings": settings})
}
