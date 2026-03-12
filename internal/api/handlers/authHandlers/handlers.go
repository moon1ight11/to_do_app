package authhandlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"todoapp/internal/api/models"
)

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	// получаем пользователя с фронта
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Name == "" || user.Pass == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}

	// добавляем пользователя в БД
	userId, err := u.userService.AddUser(user)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(userId, user.Name, user.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
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
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверяем что введенные данные не пустые
	if user.Pass == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	// поиск и проверка пользователя
	foundUser, err := u.userService.CheckAndGetUser(user)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusUnauthorized), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

// разлогин
func (u *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
