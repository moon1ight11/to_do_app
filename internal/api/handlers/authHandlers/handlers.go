package authhandlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"todoapp/internal/api/models"
)

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	// получаем пользователя с фронта
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in SignUp ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Name == "" || user.Pass == "" || user.Email == "" {
		log.Println("Error in SignUp input: name, email or password is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// добавляем пользователя в БД
	userId, err := u.userService.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Error in SignUp:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in SignUp:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(userId, user.Name, user.Email)
	if err != nil {
		log.Println("Error in Sign up:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"user_id": userId})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	// получаем пользователя с  фронта
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in Sign in ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Pass == "" || user.Email == "" {
		log.Println("Error in sign in input: email or password is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// поиск и проверка пользователя
	foundUser, err := u.userService.CheckAndGetUser(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Error in Sign in:", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		log.Println("Error in Sign in:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Email)
	if err != nil {
		log.Println("Error in Sign in:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

// разлогин
func (u *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}
