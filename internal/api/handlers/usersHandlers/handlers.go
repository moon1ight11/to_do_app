package usershandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"todoapp/internal/api/models"
)

// получение данных пользователя
func (u *UserHandler) GetUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// получаем пользователя
	user, err := u.userService.GetUser(userId)
	if err != nil {
		log.Println(err)
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
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedUser models.UserUpdate
	updatedUser.Id = userId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// обновляем нужные поля
	err := u.userService.UpdateUser(updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updatedUser": updatedUser})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// удаляем пользователя
	err := u.userService.DeleteUser(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
