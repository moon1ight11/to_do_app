package settingshandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"

	"github.com/gin-gonic/gin"
)

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.logger.Error("Error in GetSettings:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// ключ для кэша
	cacheKey := fmt.Sprintf("settings:%s", userId.String())

	// пробуем найти настройки в кэше
	if s.cacheService != nil {
		var cachedSettings models.Setting

		err := s.cacheService.Get(ctx, cacheKey, &cachedSettings)
		if err == nil {
			s.logger.Info("Settings retrieved from cache", "user", userId)
			c.JSON(http.StatusOK, gin.H{"settings": cachedSettings})
			return
		}

		// если нет - идем в БД
		s.logger.Info("Cache miss for settings", "key", cacheKey, "error", err)
	}

	// находим настройки по id
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.logger.Error("Error in GetSettings:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		s.logger.Error("Error in GetSettings:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем в кэш
	if s.cacheService != nil {
		if err := s.cacheService.Set(ctx, cacheKey, settings, 10*time.Minute); err != nil {
			s.logger.Error("Failed to set cache", "key", cacheKey, "error", err)
		}
	}

	s.logger.Info("Setting retrieved successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.logger.Error("Error in UpdateSettings:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedSettings models.UpdatedSettings

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		s.logger.Error("Error in UpdateSettings ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// получаем измененные настройки
	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		s.logger.Error("Error in UpdateSettings:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем в кэш
	if s.cacheService != nil {
		cacheKey := fmt.Sprintf("settings:%s", userId.String())
		if err := s.cacheService.Set(ctx, cacheKey, settings, 10*time.Minute); err != nil {
			s.logger.Error("Failed to update cache", "key", cacheKey, "error", err)
		}
	}

	s.logger.Info("Settings updated successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_settings": settings})
}
