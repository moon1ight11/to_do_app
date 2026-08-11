package settingshandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
	"todoapp/internal/metrics"
)

func (s *SettingsHandler) GetSettings(c *gin.Context) {
	ctx, span := s.tracer.Start(c.Request.Context(), "handler.GetSettings")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.metrics.RecordError(string(metrics.ErrForbidden), "GetSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.GetSettings: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cacheKey := fmt.Sprintf("settings:%s", userId.String())

	if s.cacheService != nil {
		var cachedSettings models.Setting
		if err := s.cacheService.Get(ctx, cacheKey, &cachedSettings); err == nil {
			s.logger.Info("settingshandlers.GetSettings: cache hit", "user", userId)
			c.JSON(http.StatusOK, gin.H{"settings": cachedSettings})
			return
		}
		s.logger.Info("settingshandlers.GetSettings: cache miss", "key", cacheKey)
	}

	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.metrics.RecordError(string(metrics.ErrInternal), "GetSettings")
			span.RecordError(err)
			s.logger.Error("settingshandlers.GetSettings: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		s.metrics.RecordError(string(metrics.ErrInternal), "GetSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.GetSettings: get settings", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if s.cacheService != nil {
		if err := s.cacheService.Set(ctx, cacheKey, settings, 10*time.Minute); err != nil {
			s.logger.Error("settingshandlers.GetSettings: cache set", "key", cacheKey, "error", err)
		}
	}

	s.logger.Info("settingshandlers.GetSettings: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	ctx, span := s.tracer.Start(c.Request.Context(), "handler.UpdateSettings")
	defer span.End()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		s.metrics.RecordError(string(metrics.ErrForbidden), "UpdateSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.UpdateSettings: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedSettings models.UpdatedSettings
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		s.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.UpdateSettings: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.settingsService.UpdateSettings(ctx, userId, updatedSettings.TimeDuration, updatedSettings.UserTz); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.metrics.RecordError(string(metrics.ErrInternal), "UpdateSettings")
			span.RecordError(err)
			s.logger.Error("settingshandlers.UpdateSettings: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		s.metrics.RecordError(string(metrics.ErrInternal), "UpdateSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.UpdateSettings: update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	settings, err := s.settingsService.GetSettings(ctx, userId)
	if err != nil {
		s.metrics.RecordError(string(metrics.ErrInternal), "UpdateSettings")
		span.RecordError(err)
		s.logger.Error("settingshandlers.UpdateSettings: get settings after update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if s.cacheService != nil {
		cacheKey := fmt.Sprintf("settings:%s", userId.String())
		if err := s.cacheService.Set(ctx, cacheKey, settings, 10*time.Minute); err != nil {
			s.logger.Error("settingshandlers.UpdateSettings: cache set", "key", cacheKey, "error", err)
		}
	}

	s.logger.Info("settingshandlers.UpdateSettings: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"updated_settings": settings})
}
