package usershandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"user": user})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if updatedUser.Name != nil {
		if strings.TrimSpace(*updatedUser.Name) == "" {
			log.Println("New name is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if updatedUser.Pass != nil {
		if strings.TrimSpace(*updatedUser.Pass) == "" {
			log.Println("New pass is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if updatedUser.Email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *updatedUser.Email)
		if err != nil {
			log.Println("Error in MatchString", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// если нет - отклоняем
		if !matched {
			log.Println("New email not looks like email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New email not looks like email"})
			return
		}
	}

	// обновляем нужные поля
	err := u.userService.UpdateUser(updatedUser.Name, updatedUser.Pass, updatedUser.Email, updatedUser.Id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": updatedUser})
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
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// проверяем существование пользователя
	_, err := u.userService.CheckUserByID(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// удаляем пользователя
	err = u.userService.DeleteUser(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// сбрасываем куки
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
