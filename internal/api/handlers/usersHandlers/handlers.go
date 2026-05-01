package usershandlers

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

// получение данных пользователя
func (u *UserHandler) GetUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.GetUser")
	defer span.End()

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in GetUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// ключ для кэша
	cacheKey := fmt.Sprintf("user:%s", userId.String())

	// пробуем найти пользователя в кэше
	if u.cacheService != nil {
		var cachedUser models.UserRequest

		err := u.cacheService.Get(ctx, cacheKey, &cachedUser)
		if err == nil {
			u.logger.Info("User retrieved from cache", "user", userId)
			c.JSON(http.StatusOK, gin.H{"user": cachedUser})
			return
		}

		// если нет - идем в БД
		u.logger.Info("Cache miss for user", "key", cacheKey, "error", err)
	}

	// получаем пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in GetUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in GetUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем в кэш
	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, user, 10*time.Minute); err != nil {
			u.logger.Error("Failed to set cache", "key", cacheKey, "error", err)
		}
	}

	u.logger.Info("User retrieved successfully", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.UpdateUser")
	defer span.End()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		span.RecordError(err)
		u.logger.Error("Error in UpdateUser ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// обновляем нужные поля
	err = u.userService.UpdateUser(ctx, updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in UpdateUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userId.String())

	// удаляем из кэша то что было до обновления
	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			u.logger.Error("Failed to invalidate user cache", "key", cacheKey, "error", err)
		}
	}

	// получаем обновленного пользователя
	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in UpdateUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// сохраняем в кэш
	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, user, 10*time.Minute); err != nil {
			u.logger.Error("Failed to set cache", "key", cacheKey, "error", err)
		}
	}

	u.logger.Info("User updated successfully", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"updated_user": user})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.DeleteUser")
	defer span.End()

	// получаем id из контекста
	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in DeleteUser:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// удаляем пользователя
	err = u.userService.DeleteUser(ctx, userId)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in DeleteUser:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in DeleteUser:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userId.String())

	// удаляем из кэша
	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			u.logger.Error("Failed to invalidate user cache", "key", cacheKey, "error", err)
		}
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User deleted successfully", "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
