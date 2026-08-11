package usershandlers

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

func (u *UserHandler) GetUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.GetUser")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrForbidden), "GetUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.GetUser: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userId.String())

	if u.cacheService != nil {
		var cachedUser models.UserRequest
		if err := u.cacheService.Get(ctx, cacheKey, &cachedUser); err == nil {
			u.logger.Info("usershandlers.GetUser: cache hit", "user", userId)
			c.JSON(http.StatusOK, gin.H{"user": cachedUser})
			return
		}
		u.logger.Info("usershandlers.GetUser: cache miss", "key", cacheKey)
	}

	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.metrics.RecordError(string(metrics.ErrInternal), "GetUser")
			span.RecordError(err)
			u.logger.Error("usershandlers.GetUser: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.metrics.RecordError(string(metrics.ErrInternal), "GetUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.GetUser: get user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, user, 10*time.Minute); err != nil {
			u.logger.Error("usershandlers.GetUser: cache set", "key", cacheKey, "error", err)
		}
	}

	u.logger.Info("usershandlers.GetUser: success", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.UpdateUser")
	defer span.End()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrForbidden), "UpdateUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.UpdateUser: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.UpdateUser: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := u.userService.UpdateUser(ctx, updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.metrics.RecordError(string(metrics.ErrInternal), "UpdateUser")
			span.RecordError(err)
			u.logger.Error("usershandlers.UpdateUser: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.metrics.RecordError(string(metrics.ErrInternal), "UpdateUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.UpdateUser: update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userId.String())

	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			u.logger.Error("usershandlers.UpdateUser: cache delete", "key", cacheKey, "error", err)
		}
	}

	user, err := u.userService.GetUser(ctx, userId)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "UpdateUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.UpdateUser: get user after update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if u.cacheService != nil {
		if err := u.cacheService.Set(ctx, cacheKey, user, 10*time.Minute); err != nil {
			u.logger.Error("usershandlers.UpdateUser: cache set", "key", cacheKey, "error", err)
		}
	}

	u.logger.Info("usershandlers.UpdateUser: success", "user", user.Id)

	c.JSON(http.StatusOK, gin.H{"updated_user": user})
}

func (u *UserHandler) DeleteUser(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.DeleteUser")
	defer span.End()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrForbidden), "DeleteUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.DeleteUser: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := u.userService.DeleteUser(ctx, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.metrics.RecordError(string(metrics.ErrInternal), "DeleteUser")
			span.RecordError(err)
			u.logger.Error("usershandlers.DeleteUser: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.metrics.RecordError(string(metrics.ErrInternal), "DeleteUser")
		span.RecordError(err)
		u.logger.Error("usershandlers.DeleteUser: delete", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userId.String())

	if u.cacheService != nil {
		if err := u.cacheService.Delete(ctx, cacheKey); err != nil {
			u.logger.Error("usershandlers.DeleteUser: cache delete", "key", cacheKey, "error", err)
		}
	}

	c.SetCookie("token", "1", -1, "/", "", false, false)

	u.logger.Info("usershandlers.DeleteUser: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
