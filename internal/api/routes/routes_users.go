package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/users"
)

type UserHandler struct {
	userService *services.UserService
}

// добаление нового пользователя в БД
func (u *UserHandler) AddUser(c *gin.Context) {
	// получаем пользователя с фронта
	var NewUser users.User
	if err := c.ShouldBindJSON(&NewUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверка валидности имени пользователя
	if strings.TrimSpace(NewUser.Name) == "" {
		log.Println("Name is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Name is empty"})
		return
	}

	// проверка валидности пароля
	if strings.TrimSpace(NewUser.Pass) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверяем, свободно ли имя пользователя
	nameExist, err := u.userService.CheckName(NewUser.Name)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	if nameExist {
		log.Println("Name already exist")
		c.JSON((http.StatusConflict), gin.H{"error": "Name already exist"})
		return
	}

	// добавляем пользователя в БД
	ID, err := u.userService.CreateUser(NewUser)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": ID})
}

// проверка существующего пользователя
func (u *UserHandler) CheckUser(c *gin.Context) {
	// получаем пользователя с фронта
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверяем пароль на валидность
	if strings.TrimSpace(user.Pass) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверка пользователя
	_, err := u.userService.CheckUser(user)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	// получаем обновленного пользователя с фронта
	var UpdatedUser struct {
		User_id    uuid.UUID `json:"user_id"`
		User_name  *string   `json:"user_name"`
		User_pass  *string   `json:"user_pass"`
		User_email *string   `json:"user_email" binding:"required,email"`
	}

	// user_id получим из куков
	UpdatedUser.User_id = uuid.New()

	if err := c.ShouldBindJSON(&UpdatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
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
	// получаем id пользователя из куков
	id := uuid.New()

	// проверяем существование пользователя
	_ , err := u.userService.CheckUserByID(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// удаляем пользователя
	err = u.userService.DeleteUser(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
