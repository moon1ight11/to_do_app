package usershandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := u.userService.GetUser(UserId)
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
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var UpdatedUser struct {
		User_id    uuid.UUID `json:"user_id"`
		User_name  *string   `json:"user_name"`
		User_pass  *string   `json:"user_pass"`
		User_email *string   `json:"user_email"`
	}

	// устанавливаем user_id
	UpdatedUser.User_id = UserId

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&UpdatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется имя - чтобы было не пустое
	if UpdatedUser.User_name != nil {
		if strings.TrimSpace(*UpdatedUser.User_name) == "" {
			log.Println("New name is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New name is empty"})
			return
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if UpdatedUser.User_pass != nil {
		if strings.TrimSpace(*UpdatedUser.User_pass) == "" {
			log.Println("New pass is empty")
			c.JSON((http.StatusBadRequest), gin.H{"error": "New pass is empty"})
			return
		}
	}

	// если обновляется почта
	if UpdatedUser.User_email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *UpdatedUser.User_email)
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
	err := u.userService.UpdateUser(UpdatedUser.User_name, UpdatedUser.User_pass, UpdatedUser.User_email, UpdatedUser.User_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": UpdatedUser})
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

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
