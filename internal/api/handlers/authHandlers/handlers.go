package authhandlers

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

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.SignUp")
	defer span.End()

	// получаем пользователя с фронта
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		span.RecordError(err)
		u.logger.Error("Error in SignUp ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Name == "" || user.Pass == "" || user.Email == "" {
		span.RecordError(fmt.Errorf("name, email and password are required"))
		u.logger.Error("Error in SignUp input: name, email or password is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// добавляем пользователя в БД
	userId, err := u.userService.AddUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in SignUp:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in SignUp:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(userId, user.Name, user.Email)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in SignUp:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	u.logger.Info("User SignUp successfully", "user", user.Email)

	c.JSON(http.StatusCreated, gin.H{"user_id": userId})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.SignIn")
	defer span.End()

	// получаем пользователя с  фронта
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		span.RecordError(err)
		u.logger.Error("Error in SignIn ShouldBindJSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Pass == "" || user.Email == "" {
		span.RecordError(fmt.Errorf("email and password are required"))
		u.logger.Error("Error in SignIn input: email or password is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// поиск и проверка пользователя
	foundUser, err := u.userService.CheckAndGetUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			u.logger.Error("Error in SignIn:", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.logger.Error("Error in SignIn:", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Email)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in SignIn:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	u.logger.Info("User SignIn successfully", "user", user.Email)

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

// разлогин
func (u *AuthHandler) SignOut(c *gin.Context) {
	_, span := u.tracer.Start(c.Request.Context(), "handler.SignOut")
	defer span.End()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		span.RecordError(err)
		u.logger.Error("Error in SignOut:", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	u.logger.Info("User SignOut successful", "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
